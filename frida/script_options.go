package frida

//#include <frida-core.h>
//#include <glib.h>
import "C"
import (
	"runtime"
	"unsafe"
)

// ScriptOptions type represents options passed to the session to create script.
type ScriptOptions struct {
	opts *C.FridaScriptOptions
}

// NewScriptOptions creates new script options with the script name provided.
func NewScriptOptions(name string) *ScriptOptions {
	opts := C.frida_script_options_new()

	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))

	C.frida_script_options_set_name(opts, nameC)

	return &ScriptOptions{
		opts: opts,
	}
}

// SetName sets the name of the script.
func (s *ScriptOptions) SetName(name string) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))

	C.frida_script_options_set_name(s.opts, nameC)
}

// SetSnapshot sets the snapshot for the script.
func (s *ScriptOptions) SetSnapshot(value []byte) {
	bts := goBytesToGBytes(value)
	runtime.SetFinalizer(bts, func(g *C.GBytes) {
		clean(unsafe.Pointer(g), unrefGObject)
	})

	C.frida_script_options_set_snapshot(s.opts, bts)
	runtime.KeepAlive(bts)
}

// SetSnapshotTransport sets the transport for the snapshot
func (s *ScriptOptions) SetSnapshotTransport(tr SnapshotOptions) {
	C.frida_script_options_set_snapshot_transport(s.opts,
		C.FridaSnapshotTransport(tr))
}

// SetRuntime sets the runtime for the script.
func (s *ScriptOptions) SetRuntime(rt ScriptRuntime) {
	C.frida_script_options_set_runtime(s.opts, C.FridaScriptRuntime(rt))
}

// GetName returns the name for the script.
func (s *ScriptOptions) GetName() string {
	return C.GoString(C.frida_script_options_get_name(s.opts))
}

// GetSnapshot returns the snapshot for the script.
func (s *ScriptOptions) GetSnapshot() []byte {
	snap := C.frida_script_options_get_snapshot(s.opts)
	bts := getGBytes(snap)
	clean(unsafe.Pointer(snap), unrefGObject)
	return bts
}

// GetSnapshotTransport returns the transport for the script.
func (s *ScriptOptions) GetSnapshotTransport() SnapshotTransport {
	tr := C.frida_script_options_get_snapshot_transport(s.opts)
	return SnapshotTransport(tr)
}
