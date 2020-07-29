package importing

import (
	"encoding/csv"
	"io"
	"log"
	"sync"
	"time"
)

// FieldMapping - field mapping
type FieldMapping map[string]int

// Importer - Importer object used to parse csv and save rows to DB
type Importer struct {
	Timestamp         time.Time
	Comma             rune
	Comment           rune
	NumberOfChannels  int
	ChannelBufferSize int
	channels          []chan rowToSave
	wg                *sync.WaitGroup
	Mapping           FieldMapping
	ParseRow          func([]string, FieldMapping) (interface{}, error)
	EntitySaver       func(interface{}) error
}

type rowToSave struct {
	timestamp time.Time
	entity    interface{}
}

// Parse - parse provided payload
func (im Importer) Parse(reader io.Reader) {
	im.createChannels()
	defer im.closeChannels()
	im.runGoRoutines()
	r := csv.NewReader(reader)
	fields, err := r.Read()
	if err != nil {
		log.Printf("Importing error at time %d", im.Timestamp.Unix())
	}
	for i := 0; err == nil && fields != nil; i++ {
		entity, err := im.ParseRow(fields, im.Mapping)
		if err != nil {
			log.Printf("Importing error at time %d", im.Timestamp.Unix())
		}
		im.channels[i%im.NumberOfChannels] <- rowToSave{
			timestamp: im.Timestamp,
			entity:    entity,
		}
		fields, err = r.Read()
	}
	im.wg.Wait()
}

func (im Importer) createChannels() {
	im.channels = make([]chan rowToSave, im.NumberOfChannels)
	for i := range im.channels {
		im.channels[i] = make(chan rowToSave, im.ChannelBufferSize)
	}
}

func (im Importer) closeChannels() {
	for _, c := range im.channels {
		close(c)
	}
}

func (im Importer) runGoRoutines() {
	im.wg = &sync.WaitGroup{}
	im.wg.Add(im.NumberOfChannels)
	for _, c := range im.channels {
		go entitySaverRoutine(c, im.wg, im.EntitySaver)
	}
}

func entitySaverRoutine(c <-chan rowToSave, wg *sync.WaitGroup, saveEntity func(interface{}) error) {
	var err error
	for i := range c {
		if err = saveEntity(i); err != nil {
			log.Printf("[Import error - %d] Could not save entity %v", i.timestamp.Unix(), i.entity)
		}
	}
	wg.Done()
}
