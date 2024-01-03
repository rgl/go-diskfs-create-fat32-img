package main

import (
	"fmt"
	"log"
	"os"

	diskfs "github.com/diskfs/go-diskfs"
	diskpkg "github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/partition/mbr"
)

func main() {
	amtSetup, err := os.ReadFile("Setup.bin")
	if err != nil {
		log.Panic(err)
	}

	// NB AMT Setup.bin must be stored inside a FAT32 filesystem.

	sectorSize := int(diskfs.SectorSize512)
	// align the partition to 1 MiB.
	partitionStartSector := (1 * 1024 * 1024) / sectorSize
	// NB the FAT32 minimum size is 33 MiB.
	partitionSizeSector := max(33*1024*1024, len(amtSetup)) / sectorSize
	diskSize := (partitionStartSector + partitionSizeSector) * sectorSize

	partitionTable := &mbr.Table{
		LogicalSectorSize:  sectorSize,
		PhysicalSectorSize: sectorSize,
		Partitions: []*mbr.Partition{
			{
				Type:  mbr.Fat32CHS,
				Start: uint32(partitionStartSector),
				Size:  uint32(partitionSizeSector),
			},
		},
	}

	disk, err := diskfs.Create("Setup.bin.img", int64(diskSize), diskfs.Raw, diskfs.SectorSize(sectorSize))
	if err != nil {
		log.Panic(err)
	}

	err = disk.Partition(partitionTable)
	if err != nil {
		log.Panic(err)
	}

	spec := diskpkg.FilesystemSpec{
		Partition:   1,
		FSType:      filesystem.TypeFat32,
		VolumeLabel: "AMT-SETUP", // NB this can be at most 11 chars.
	}
	fs, err := disk.CreateFilesystem(spec)
	if err != nil {
		log.Panic(err)
	}

	amtSetupFile, err := fs.OpenFile("/Setup.bin", os.O_CREATE|os.O_RDWR)
	if err != nil {
		log.Panic(err)
	}

	_, err = amtSetupFile.Write(amtSetup)
	if err != nil {
		log.Panic(err)
	}

	fmt.Print(`
Setup.bin.img can be mounted as:

fdisk -l Setup.bin.img --bytes
sudo mkdir /mnt/amt-setup
target_device="$(sudo losetup --partscan --show --find Setup.bin.img)"
sudo fdisk -l "$target_device"
sudo mount -o ro "${target_device}p1" /mnt/amt-setup
sudo df -h /mnt/amt-setup
sudo ls -laF /mnt/amt-setup
sudo sha256sum /mnt/amt-setup/Setup.bin Setup.bin
sudo umount /mnt/amt-setup
sudo losetup --detach "$target_device"
`)
}
