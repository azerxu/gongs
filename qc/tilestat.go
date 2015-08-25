package qc

import (
	"fmt"
	"gongs/biofile/fastq"
	"gongs/stat"
	"gongs/xopen"
	"sort"
	"strconv"
	"strings"
)

type flowcell struct {
	id     string
	length int
	lanes  map[int]*lane
}

type lane struct {
	id    int
	tiles map[int]*tile
}

type tile struct {
	id     int
	cycles map[int]*cycle
}

type cycle struct {
	pos int
	m   *stat.IntMap
}

func (c *cycle) meanQ() float64 {
	return c.m.Mean()
}

func (c *cycle) medianQ() float64 {
	return c.m.Median()
}

func (c *cycle) maxQ() int {
	maxq := 0
	for q := range c.m.Data {
		if maxq < q {
			maxq = q
		}
	}
	return maxq
}

func (c *cycle) minQ() int {
	minq := 127
	for q := range c.m.Data {
		if minq > q {
			minq = q
		}
	}
	return minq
}

func (c *cycle) percentile10() float64 {
	return c.m.Percentile(0.10)
}

func (c *cycle) percentile25() float64 {
	return c.m.Percentile(0.25)
}

func (c *cycle) percentile75() float64 {
	return c.m.Percentile(0.75)
}

func (c *cycle) percentile90() float64 {
	return c.m.Percentile(0.90)
}

type Tilestat struct {
	flowcells   map[string]*flowcell // record quality by flowcell,lane,tile
	quals       map[int]int          // record full quality count distribution
	qualByCycle map[int]*cycle       // record quality by position
	max         byte                 // max quality score
	min         byte                 // min quality score
}

func NewTile() *Tilestat {
	return &Tilestat{
		flowcells:   make(map[string]*flowcell),
		quals:       make(map[int]int),
		qualByCycle: make(map[int]*cycle),
		max:         0,
		min:         127,
	}
}

// Count count quality by postion and flowcell,lane,tile
func (t *Tilestat) Count(fq *fastq.Fastq) error {
	ids := strings.SplitN(fq.Name, ":", 6)
	// fmt.Println(ids)
	flowid := ids[2]
	laneid, err := strconv.Atoi(ids[3])
	if err != nil {
		return err
	}
	tileid, err := strconv.Atoi(ids[4])
	if err != nil {
		return err
	}

	mflowcell, ok := t.flowcells[flowid]
	if !ok {
		mflowcell = &flowcell{id: flowid, lanes: make(map[int]*lane)}
		t.flowcells[flowid] = mflowcell
	}
	if mflowcell.length < len(fq.Seq) {
		mflowcell.length = len(fq.Seq)
	}

	mlane, ok := mflowcell.lanes[laneid]
	if !ok {
		mlane = &lane{id: laneid, tiles: make(map[int]*tile)}
		mflowcell.lanes[laneid] = mlane
	}
	mtile, ok := mlane.tiles[tileid]
	if !ok {
		mtile = &tile{id: tileid, cycles: make(map[int]*cycle)}
		mlane.tiles[tileid] = mtile
	}

	for i, q := range fq.Qual { // record each quality
		t.quals[int(q)]++
		ncycle, ok := t.qualByCycle[i]
		if !ok {
			ncycle = &cycle{pos: i, m: stat.NewIntMap(make(map[int]int))}
			t.qualByCycle[i] = ncycle
		}
		ncycle.m.Data[int(q)]++

		if t.min > q {
			t.min = q
		}
		if t.max < q {
			t.max = q
		}
		mcycle, ok := mtile.cycles[i]
		if !ok {
			mcycle = &cycle{pos: i, m: stat.NewIntMap(map[int]int{})}
			mtile.cycles[i] = mcycle
		}
		mcycle.m.Data[int(q)]++
	}
	return nil
}

func (t *Tilestat) MinQual() int {
	return int(t.min)
}

func (t *Tilestat) MaxQual() int {
	return int(t.max)
}

// GuessEncoding Guess quality Encoding version
// S - Sanger        Phred+33,  raw reads typically (0, 40), using ASCII 33 to 73
// X - Solexa        Solexa+64, raw reads typically (-5, 40), using ASCII 59 to 104
// I - Illumina 1.3+ Phred+64,  raw reads typically (0, 40), using ASCII 64 to 104
// J - Illumina 1.5+ Phred+64,  raw reads typically (3, 40), using ASCII 66 to 104
// L - Illumina 1.8+ Phred+33,  raw reads typically (0, 41), using ASCII 33 to 74
func (t *Tilestat) GuessEncoding() string {
	if t.min < 33 {
		return "Unkown"
	} else if t.max < 74 {
		return "SANGER"
	} else if t.max < 75 {
		return "Illumina1.8+"
	} else if t.min > 58 && t.min < 64 {
		return "Solexa"
	} else if t.min > 63 && t.min < 66 {
		return "Illumina1.3+"
	} else if t.min > 65 {
		return "Illumina11.5+"
	}
	// minqual < 58 and maxqual > 74
	return "Mix"
}

func (t *Tilestat) Q20() float64 {
	return t.Q(20)
}

func (t *Tilestat) Q30() float64 {
	return t.Q(30)
}

// Q return >=qual percent
func (t *Tilestat) Q(q byte) float64 {
	c := 0
	tot := 0
	for qu, count := range t.quals {
		tot += count
		if qu < int(q) {
			continue
		}
		c += count
	}
	return float64(c*100) / float64(tot)
}

// Qat return qual count
func (t *Tilestat) Qat(q byte) int {
	return t.quals[int(q)]
}

// Save file block
// General File Format Option:
//  1. #! line is comment line
//  2. #: line is key value line, eg. #: flowcell: abcdefg
//  3. ## line is header line

// SaveQualDist save whole quality distribution
func (t *Tilestat) SaveQualDist(prefix string) error {
	f, err := xopen.Xcreate(prefix+".qualdist", "w")
	if err != nil {
		return err
	}
	defer f.Close()

	// print header line
	header := []string{"MinQ", " MaxQ"}
	for q := t.min; q < t.max+1; q++ {
		header = append(header, fmt.Sprintf("Q%d", q))
	}
	fmt.Fprintln(f, "##", strings.Join(header, "\t"))

	// print data line
	fmt.Fprint(f, strconv.Itoa(int(t.min)), "\t", strconv.Itoa(int(t.max)))
	for i, max := int(t.min), int(t.max); i < max+1; i++ {
		fmt.Fprint(f, "\t", strconv.Itoa(t.quals[i]))
	}
	return nil
}

func (t *Tilestat) SaveCycleStat(prefix string) error {
	f, err := xopen.Xcreate(prefix+".cyclestat", "w")
	if err != nil {
		return err
	}
	defer f.Close()

	// print header line
	header := []string{"cycle", "minQ", "maxQ", "meanQ", "medianQ", "percentile10", "percentile25", "percentile75", "percentile90"}
	fmt.Fprintln(f, "##", strings.Join(header, "\t"))

	// print data line
	for i, l := 0, len(t.qualByCycle); i < l; i++ {
		c := t.qualByCycle[i]
		fmt.Fprintln(f, i, "\t", c.minQ(), "\t", c.maxQ(), "\t", c.meanQ(), "\t", c.medianQ(), "\t",
			c.percentile10(), "\t", c.percentile25(), "\t", c.percentile75(), "\t", c.percentile90())
	}
	return nil
}

// SaveTileStat save Tile qualtity stat
func (t *Tilestat) SaveTileStat(prefix string) error {
	f, err := xopen.Xcreate(prefix+".tilestat", "w")
	if err != nil {
		return err
	}
	defer f.Close()

	for flowid := range t.flowcells {
		mflow := t.flowcells[flowid]
		fmt.Fprintln(f, "#!", strings.Repeat("=", 20), flowid, strings.Repeat("=", 20))

		// print header line
		header := []string{"flowid", "laneid", "tileid", "length"}
		for i := 0; i < mflow.length; i++ {
			header = append(header, fmt.Sprintf("%d", i+1))
		}
		fmt.Fprintln(f, "##", strings.Join(header, "\t"))

		// get sorted  laneids
		laneids := []int{}
		for laneid := range mflow.lanes {
			laneids = append(laneids, laneid)
		}
		sort.Ints(laneids)

		for _, laneid := range laneids { // iterate each lane
			mlane := mflow.lanes[laneid]

			// get sorted tileids
			tileids := []int{}
			for tileid := range mlane.tiles {
				tileids = append(tileids, tileid)
			}
			sort.Ints(tileids)

			for _, tileid := range tileids { // iterate each tile
				mtile := mlane.tiles[tileid]
				result := []string{flowid, strconv.Itoa(laneid), strconv.Itoa(tileid), strconv.Itoa(mflow.length)}
				for i := 0; i < mflow.length; i++ {
					cycle, ok := mtile.cycles[i]
					if !ok {
						result = append(result, "0")
						continue
					}
					result = append(result, fmt.Sprintf("%.1f", cycle.medianQ()))
				}
				// print tile lane
				fmt.Fprintln(f, strings.Join(result, "\t"))
			}
		}
	}
	return nil
}
