package managedplugin

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

// ErrArchMismatch is returned when a cached plugin binary was built for a
// different architecture than the one we are running on.
var ErrArchMismatch = errors.New("plugin binary architecture mismatch")

var elfMachines = map[string]elf.Machine{
	"386":      elf.EM_386,
	"amd64":    elf.EM_X86_64,
	"arm":      elf.EM_ARM,
	"arm64":    elf.EM_AARCH64,
	"loong64":  elf.EM_LOONGARCH,
	"mips":     elf.EM_MIPS,
	"mips64":   elf.EM_MIPS,
	"mips64le": elf.EM_MIPS,
	"mipsle":   elf.EM_MIPS,
	"ppc64":    elf.EM_PPC64,
	"ppc64le":  elf.EM_PPC64,
	"riscv64":  elf.EM_RISCV,
	"s390x":    elf.EM_S390,
}

var machoCPUs = map[string]macho.Cpu{
	"386":   macho.Cpu386,
	"amd64": macho.CpuAmd64,
	"arm":   macho.CpuArm,
	"arm64": macho.CpuArm64,
}

var peMachines = map[string]uint16{
	"386":   pe.IMAGE_FILE_MACHINE_I386,
	"amd64": pe.IMAGE_FILE_MACHINE_AMD64,
	"arm":   pe.IMAGE_FILE_MACHINE_ARMNT,
	"arm64": pe.IMAGE_FILE_MACHINE_ARM64,
}

// pluginCachePaths returns the paths a plugin binary may be cached at. The
// canonical path is namespaced by GOOS_GOARCH so that hosts of different
// architectures sharing a cq directory cannot overwrite each other. The legacy
// path is the pre-namespacing location, which is read but never written.
func pluginCachePaths(directory, kind, org, name, version string) (canonical, legacy string) {
	base := filepath.Join(directory, "plugins", kind, org, name, version)
	target := runtime.GOOS + "_" + runtime.GOARCH
	canonical = WithBinarySuffix(filepath.Join(base, target, "plugin"))
	legacy = WithBinarySuffix(filepath.Join(base, "plugin"))
	return canonical, legacy
}

// resolveCachedPlugin reports which of the cache locations holds a binary that
// is usable on this host, preferring the canonical one. A legacy binary is
// reused only when it was built for the architecture we are running on, so
// existing caches survive without risking an exec format error.
func resolveCachedPlugin(logger zerolog.Logger, canonical, legacy string) (string, bool) {
	if err := validateBinary(canonical); err == nil {
		return canonical, true
	} else if !os.IsNotExist(err) {
		logger.Warn().Str("path", canonical).Err(err).Msg("cached plugin is unusable, re-downloading")
		return canonical, false
	}

	if err := validateBinary(legacy); err == nil {
		return legacy, true
	} else if !os.IsNotExist(err) {
		logger.Warn().Str("path", legacy).Err(err).Msg("cached plugin at legacy path is unusable, downloading to architecture specific path")
	}

	return canonical, false
}

// validateBinary checks that path holds a complete executable that can run on
// this host. It guards against two ways a shared cache directory goes bad: a
// binary left behind by a host of a different architecture, and a binary
// truncated by an interrupted or concurrent download.
func validateBinary(path string) error {
	switch runtime.GOOS {
	case "darwin", "ios":
		return validateMachO(path)
	case "windows":
		return validatePE(path)
	default:
		return validateELF(path)
	}
}

func validateELF(path string) error {
	f, err := elf.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	want, ok := elfMachines[runtime.GOARCH]
	if !ok {
		return nil
	}
	if f.Machine != want {
		return fmt.Errorf("%w: binary is %s, expected %s", ErrArchMismatch, f.Machine, want)
	}
	return nil
}

func validateMachO(path string) error {
	f, err := macho.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	want, ok := machoCPUs[runtime.GOARCH]
	if !ok {
		return nil
	}
	if f.Cpu != want {
		return fmt.Errorf("%w: binary is %s, expected %s", ErrArchMismatch, f.Cpu, want)
	}
	return nil
}

func validatePE(path string) error {
	f, err := pe.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	want, ok := peMachines[runtime.GOARCH]
	if !ok {
		return nil
	}
	if f.Machine != want {
		return fmt.Errorf("%w: binary machine is %#x, expected %#x", ErrArchMismatch, f.Machine, want)
	}
	return nil
}
