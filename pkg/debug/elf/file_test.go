// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package elf

import (
	"debug/dwarf"
	"encoding/binary"
	"reflect"
	"testing"
)

type fileTest struct {
	file     string
	hdr      FileHeader
	sections []SectionHeader
}

var fileTests = []fileTest{
	fileTest{
		"testdata/gcc-386-freebsd-exec",
		FileHeader{ELFCLASS32, ELFDATA2LSB, EV_CURRENT, ELFOSABI_FREEBSD, 0, binary.LittleEndian, ET_EXEC, EM_386},
		[]SectionHeader{
			SectionHeader{"", SHT_NULL, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			SectionHeader{".interp", SHT_PROGBITS, SHF_ALLOC, 0x80480d4, 0xd4, 0x15, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".hash", SHT_HASH, SHF_ALLOC, 0x80480ec, 0xec, 0x90, 0x3, 0x0, 0x4, 0x4},
			SectionHeader{".dynsym", SHT_DYNSYM, SHF_ALLOC, 0x804817c, 0x17c, 0x110, 0x4, 0x1, 0x4, 0x10},
			SectionHeader{".dynstr", SHT_STRTAB, SHF_ALLOC, 0x804828c, 0x28c, 0xbb, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".rel.plt", SHT_REL, SHF_ALLOC, 0x8048348, 0x348, 0x20, 0x3, 0x7, 0x4, 0x8},
			SectionHeader{".init", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x8048368, 0x368, 0x11, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".plt", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x804837c, 0x37c, 0x50, 0x0, 0x0, 0x4, 0x4},
			SectionHeader{".text", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x80483cc, 0x3cc, 0x180, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".fini", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x804854c, 0x54c, 0xc, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".rodata", SHT_PROGBITS, SHF_ALLOC, 0x8048558, 0x558, 0xa3, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".data", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x80495fc, 0x5fc, 0xc, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".eh_frame", SHT_PROGBITS, SHF_ALLOC, 0x8049608, 0x608, 0x4, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".dynamic", SHT_DYNAMIC, SHF_WRITE + SHF_ALLOC, 0x804960c, 0x60c, 0x98, 0x4, 0x0, 0x4, 0x8},
			SectionHeader{".ctors", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x80496a4, 0x6a4, 0x8, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".dtors", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x80496ac, 0x6ac, 0x8, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".jcr", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x80496b4, 0x6b4, 0x4, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".got", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x80496b8, 0x6b8, 0x1c, 0x0, 0x0, 0x4, 0x4},
			SectionHeader{".bss", SHT_NOBITS, SHF_WRITE + SHF_ALLOC, 0x80496d4, 0x6d4, 0x20, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".comment", SHT_PROGBITS, 0x0, 0x0, 0x6d4, 0x12d, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_aranges", SHT_PROGBITS, 0x0, 0x0, 0x801, 0x20, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_pubnames", SHT_PROGBITS, 0x0, 0x0, 0x821, 0x1b, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_info", SHT_PROGBITS, 0x0, 0x0, 0x83c, 0x11d, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_abbrev", SHT_PROGBITS, 0x0, 0x0, 0x959, 0x41, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_line", SHT_PROGBITS, 0x0, 0x0, 0x99a, 0x35, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_frame", SHT_PROGBITS, 0x0, 0x0, 0x9d0, 0x30, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".debug_str", SHT_PROGBITS, 0x0, 0x0, 0xa00, 0xd, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".shstrtab", SHT_STRTAB, 0x0, 0x0, 0xa0d, 0xf8, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".symtab", SHT_SYMTAB, 0x0, 0x0, 0xfb8, 0x4b0, 0x1d, 0x38, 0x4, 0x10},
			SectionHeader{".strtab", SHT_STRTAB, 0x0, 0x0, 0x1468, 0x206, 0x0, 0x0, 0x1, 0x0},
		},
	},
	fileTest{
		"testdata/gcc-amd64-linux-exec",
		FileHeader{ELFCLASS64, ELFDATA2LSB, EV_CURRENT, ELFOSABI_NONE, 0, binary.LittleEndian, ET_EXEC, EM_X86_64},
		[]SectionHeader{
			SectionHeader{"", SHT_NULL, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			SectionHeader{".interp", SHT_PROGBITS, SHF_ALLOC, 0x400200, 0x200, 0x1c, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".note.ABI-tag", SHT_NOTE, SHF_ALLOC, 0x40021c, 0x21c, 0x20, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".hash", SHT_HASH, SHF_ALLOC, 0x400240, 0x240, 0x24, 0x5, 0x0, 0x8, 0x4},
			SectionHeader{".gnu.hash", SHT_LOOS + 268435446, SHF_ALLOC, 0x400268, 0x268, 0x1c, 0x5, 0x0, 0x8, 0x0},
			SectionHeader{".dynsym", SHT_DYNSYM, SHF_ALLOC, 0x400288, 0x288, 0x60, 0x6, 0x1, 0x8, 0x18},
			SectionHeader{".dynstr", SHT_STRTAB, SHF_ALLOC, 0x4002e8, 0x2e8, 0x3d, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".gnu.version", SHT_HIOS, SHF_ALLOC, 0x400326, 0x326, 0x8, 0x5, 0x0, 0x2, 0x2},
			SectionHeader{".gnu.version_r", SHT_LOOS + 268435454, SHF_ALLOC, 0x400330, 0x330, 0x20, 0x6, 0x1, 0x8, 0x0},
			SectionHeader{".rela.dyn", SHT_RELA, SHF_ALLOC, 0x400350, 0x350, 0x18, 0x5, 0x0, 0x8, 0x18},
			SectionHeader{".rela.plt", SHT_RELA, SHF_ALLOC, 0x400368, 0x368, 0x30, 0x5, 0xc, 0x8, 0x18},
			SectionHeader{".init", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x400398, 0x398, 0x18, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".plt", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x4003b0, 0x3b0, 0x30, 0x0, 0x0, 0x4, 0x10},
			SectionHeader{".text", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x4003e0, 0x3e0, 0x1b4, 0x0, 0x0, 0x10, 0x0},
			SectionHeader{".fini", SHT_PROGBITS, SHF_ALLOC + SHF_EXECINSTR, 0x400594, 0x594, 0xe, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".rodata", SHT_PROGBITS, SHF_ALLOC, 0x4005a4, 0x5a4, 0x11, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".eh_frame_hdr", SHT_PROGBITS, SHF_ALLOC, 0x4005b8, 0x5b8, 0x24, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".eh_frame", SHT_PROGBITS, SHF_ALLOC, 0x4005e0, 0x5e0, 0xa4, 0x0, 0x0, 0x8, 0x0},
			SectionHeader{".ctors", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x600688, 0x688, 0x10, 0x0, 0x0, 0x8, 0x0},
			SectionHeader{".dtors", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x600698, 0x698, 0x10, 0x0, 0x0, 0x8, 0x0},
			SectionHeader{".jcr", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x6006a8, 0x6a8, 0x8, 0x0, 0x0, 0x8, 0x0},
			SectionHeader{".dynamic", SHT_DYNAMIC, SHF_WRITE + SHF_ALLOC, 0x6006b0, 0x6b0, 0x1a0, 0x6, 0x0, 0x8, 0x10},
			SectionHeader{".got", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x600850, 0x850, 0x8, 0x0, 0x0, 0x8, 0x8},
			SectionHeader{".got.plt", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x600858, 0x858, 0x28, 0x0, 0x0, 0x8, 0x8},
			SectionHeader{".data", SHT_PROGBITS, SHF_WRITE + SHF_ALLOC, 0x600880, 0x880, 0x18, 0x0, 0x0, 0x8, 0x0},
			SectionHeader{".bss", SHT_NOBITS, SHF_WRITE + SHF_ALLOC, 0x600898, 0x898, 0x8, 0x0, 0x0, 0x4, 0x0},
			SectionHeader{".comment", SHT_PROGBITS, 0x0, 0x0, 0x898, 0x126, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_aranges", SHT_PROGBITS, 0x0, 0x0, 0x9c0, 0x90, 0x0, 0x0, 0x10, 0x0},
			SectionHeader{".debug_pubnames", SHT_PROGBITS, 0x0, 0x0, 0xa50, 0x25, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_info", SHT_PROGBITS, 0x0, 0x0, 0xa75, 0x1a7, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_abbrev", SHT_PROGBITS, 0x0, 0x0, 0xc1c, 0x6f, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_line", SHT_PROGBITS, 0x0, 0x0, 0xc8b, 0x13f, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".debug_str", SHT_PROGBITS, SHF_MERGE + SHF_STRINGS, 0x0, 0xdca, 0xb1, 0x0, 0x0, 0x1, 0x1},
			SectionHeader{".debug_ranges", SHT_PROGBITS, 0x0, 0x0, 0xe80, 0x90, 0x0, 0x0, 0x10, 0x0},
			SectionHeader{".shstrtab", SHT_STRTAB, 0x0, 0x0, 0xf10, 0x149, 0x0, 0x0, 0x1, 0x0},
			SectionHeader{".symtab", SHT_SYMTAB, 0x0, 0x0, 0x19a0, 0x6f0, 0x24, 0x39, 0x8, 0x18},
			SectionHeader{".strtab", SHT_STRTAB, 0x0, 0x0, 0x2090, 0x1fc, 0x0, 0x0, 0x1, 0x0},
		},
	},
}

func TestOpen(t *testing.T) {
	for i := range fileTests {
		tt := &fileTests[i]

		f, err := Open(tt.file)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(f.FileHeader, tt.hdr) {
			t.Errorf("open %s:\n\thave %#v\n\twant %#v\n", tt.file, f.FileHeader, tt.hdr)
			continue
		}
		for i, s := range f.Sections {
			if i >= len(tt.sections) {
				break
			}
			sh := &tt.sections[i]
			if !reflect.DeepEqual(&s.SectionHeader, sh) {
				t.Errorf("open %s, section %d:\n\thave %#v\n\twant %#v\n", tt.file, i, &s.SectionHeader, sh)
			}
		}
		tn := len(tt.sections)
		fn := len(f.Sections)
		if tn != fn {
			t.Errorf("open %s: len(Sections) = %d, want %d", tt.file, fn, tn)
		}
	}
}

type relocationTest struct {
	file       string
	firstEntry *dwarf.Entry
}

var relocationTests = []relocationTest{
	relocationTest{
		"testdata/go-relocation-test-gcc441-x86-64.o",
		&dwarf.Entry{Offset: 0xb, Tag: dwarf.TagCompileUnit, Children: true, Field: []dwarf.Field{dwarf.Field{Attr: dwarf.AttrProducer, Val: "GNU C 4.4.1"}, dwarf.Field{Attr: dwarf.AttrLanguage, Val: int64(1)}, dwarf.Field{Attr: dwarf.AttrName, Val: "go-relocation-test.c"}, dwarf.Field{Attr: dwarf.AttrCompDir, Val: "/tmp"}, dwarf.Field{Attr: dwarf.AttrLowpc, Val: uint64(0x0)}, dwarf.Field{Attr: dwarf.AttrHighpc, Val: uint64(0x6)}, dwarf.Field{Attr: dwarf.AttrStmtList, Val: int64(0)}}},
	},
	relocationTest{
		"testdata/go-relocation-test-gcc441-x86.o",
		&dwarf.Entry{Offset: 0xb, Tag: dwarf.TagCompileUnit, Children: true, Field: []dwarf.Field{dwarf.Field{Attr: dwarf.AttrProducer, Val: "GNU C 4.4.1"}, dwarf.Field{Attr: dwarf.AttrLanguage, Val: int64(1)}, dwarf.Field{Attr: dwarf.AttrName, Val: "t.c"}, dwarf.Field{Attr: dwarf.AttrCompDir, Val: "/tmp"}, dwarf.Field{Attr: dwarf.AttrLowpc, Val: uint64(0x0)}, dwarf.Field{Attr: dwarf.AttrHighpc, Val: uint64(0x5)}, dwarf.Field{Attr: dwarf.AttrStmtList, Val: int64(0)}}},
	},
	relocationTest{
		"testdata/go-relocation-test-gcc424-x86-64.o",
		&dwarf.Entry{Offset: 0xb, Tag: dwarf.TagCompileUnit, Children: true, Field: []dwarf.Field{dwarf.Field{Attr: dwarf.AttrProducer, Val: "GNU C 4.2.4 (Ubuntu 4.2.4-1ubuntu4)"}, dwarf.Field{Attr: dwarf.AttrLanguage, Val: int64(1)}, dwarf.Field{Attr: dwarf.AttrName, Val: "go-relocation-test-gcc424.c"}, dwarf.Field{Attr: dwarf.AttrCompDir, Val: "/tmp"}, dwarf.Field{Attr: dwarf.AttrLowpc, Val: uint64(0x0)}, dwarf.Field{Attr: dwarf.AttrHighpc, Val: uint64(0x6)}, dwarf.Field{Attr: dwarf.AttrStmtList, Val: int64(0)}}},
	},
}

func TestDWARFRelocations(t *testing.T) {
	for i, test := range relocationTests {
		f, err := Open(test.file)
		if err != nil {
			t.Error(err)
			continue
		}
		dwarf, err := f.DWARF()
		if err != nil {
			t.Error(err)
			continue
		}
		reader := dwarf.Reader()
		// Checking only the first entry is sufficient since it has
		// many different strings. If the relocation had failed, all
		// the string offsets would be zero and all the strings would
		// end up being the same.
		firstEntry, err := reader.Next()
		if err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(test.firstEntry, firstEntry) {
			t.Errorf("#%d: mismatch: got:%#v want:%#v", i, firstEntry, test.firstEntry)
			continue
		}
	}
}
