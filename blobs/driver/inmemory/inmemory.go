package inmemory

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/danielkrainas/gobag/decouple/drivers"
	"github.com/danielkrainas/gobag/util/bagio"

	"github.com/danielkrainas/shexd/blobs"
	"github.com/danielkrainas/shexd/blobs/driver/factory"
)

type driverFactory struct{}

func (df *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	return &driver{
		blobs: make(map[string]*blobDescriptor, 0),
	}, nil
}

func init() {
	factory.Register("inmemory", &driverFactory{})
}

type blobDescriptor struct {
	blob *blobs.Blob
	data *bytes.Buffer
}

type driver struct {
	m     sync.Mutex
	blobs map[string]*blobDescriptor
}

func (d *driver) Inspect(name string) (*blobs.Blob, error) {
	d.m.Lock()
	defer d.m.Unlock()

	desc, ok := d.blobs[name]
	if !ok {
		return nil, blobs.ErrUnknown
	}

	return desc.blob, nil
}

func (d *driver) WriteMeta(name string, b *blobs.Blob) error {
	d.m.Lock()
	defer d.m.Unlock()

	desc, ok := d.blobs[name]
	if !ok {
		return blobs.ErrUnknown
	}

	desc.blob = b
	if name != b.Name {
		delete(d.blobs, name)
		d.blobs[b.Name] = desc
	}

	return nil
}

func (d *driver) Writer(name string) (io.WriteCloser, error) {
	d.m.Lock()
	defer d.m.Unlock()

	desc, ok := d.blobs[name]
	if !ok {
		return nil, blobs.ErrUnknown
	}

	w := bytes.NewBuffer(make([]byte, 0, 0))
	desc.data = w
	return bagio.NopWriteCloser(w), nil
}

func (d *driver) Reader(name string) (io.ReadCloser, error) {
	d.m.Lock()
	defer d.m.Unlock()

	desc, ok := d.blobs[name]
	if !ok {
		return nil, blobs.ErrUnknown
	}

	b := desc.data.Bytes()[:]
	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func (d *driver) Drop(name string) (bool, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.blobs[name]; !ok {
		return false, nil
	}

	delete(d.blobs, name)
	return true, nil
}
