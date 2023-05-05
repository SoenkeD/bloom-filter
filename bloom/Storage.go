package bloom

import (
	"log"
	"strconv"
)

type StorageInitInput struct {
	MinEntryCount uint64
}

type Storage struct {
	rawStorage []uint64
}

func InitStorage(input StorageInitInput) (storage *Storage) {

	requiredInts := input.MinEntryCount/64 + 1
	rawStorage := make([]uint64, requiredInts)

	storage = &Storage{
		rawStorage: rawStorage,
	}

	return
}

func (storage *Storage) getIdx(hash uint64) (relIntIdx int64, relBitIdx int64, relInt uint64) {

	convertedHash := int64(hash)
	if convertedHash < 0 {
		log.Fatal("invalid hash received", convertedHash, hash)
	}

	relIntIdx = convertedHash / 64
	relBitIdx = convertedHash % 64
	relInt = storage.rawStorage[relIntIdx]
	return
}

func (storage *Storage) GetSlotCount() (slotCount int64) {
	slotCount = int64(len(storage.rawStorage) * 64)
	return
}

func (storage *Storage) AddItem(hash uint64) (err error) {
	relIntIdx, relBitIdx, relInt := storage.getIdx(hash)
	convertedInt := getConvertedInt(relInt)

	prefix := ""
	if relBitIdx != 0 {
		prefix = convertedInt[:relBitIdx]
	}

	postfix := ""
	if relBitIdx != 63 {
		postfix = convertedInt[relBitIdx+1:]
	}

	newInt := prefix + "1" + postfix
	parsedUint, err := strconv.ParseUint(newInt, 2, 64)
	if err != nil {
		return
	}

	storage.rawStorage[relIntIdx] = parsedUint

	return

}

func (storage *Storage) PotentiallyKnowItem(hash uint64) bool {

	_, relBitIdx, relInt := storage.getIdx(hash)
	convertedInt := getConvertedInt(relInt)

	bitVal := convertedInt[relBitIdx : relBitIdx+1]

	if bitVal == "1" {
		return true
	}

	return false
}

func getConvertedInt(toFillInt uint64) (filledString string) {

	filledString = strconv.FormatUint(toFillInt, 2)
	for len(filledString) < 64 {
		filledString = "0" + filledString
	}

	return
}
