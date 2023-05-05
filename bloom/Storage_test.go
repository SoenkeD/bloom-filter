package bloom

import (
	"testing"
)

func storageBeforeEach() (input StorageInitInput, storage *Storage) {
	input = StorageInitInput{MinEntryCount: 2 * 64}
	storage = InitStorage(input)
	return
}

func TestInitStorage(t *testing.T) {

	// default before each -> one overshoot (keep in mind 0 is included)
	_, storage := storageBeforeEach()
	if len(storage.rawStorage) != 3 {
		t.Errorf("expected=2 got=%d", len(storage.rawStorage))
	}

	// fits
	input := StorageInitInput{MinEntryCount: 2*64 - 1}
	storage = InitStorage(input)
	if len(storage.rawStorage) != 2 {
		t.Errorf("expected=2 got=%d", len(storage.rawStorage))
	}
}

func TestStorage_getIdx(t *testing.T) {
	_, storage := storageBeforeEach()

	// very first bit
	relIntIdx, relBitIdx, _ := storage.getIdx(0)
	if relIntIdx != 0 {
		t.Errorf("expected=0 got=%d", relIntIdx)
	}
	if relBitIdx != 0 {
		t.Errorf("expected=0 got=%d", relBitIdx)
	}

	// 2nd bit in first chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(1)
	if relIntIdx != 0 {
		t.Errorf("expected=0 got=%d", relIntIdx)
	}
	if relBitIdx != 1 {
		t.Errorf("expected=1 got=%d", relBitIdx)
	}

	// last bit in first chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(63)
	if relIntIdx != 0 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 63 {
		t.Errorf("expected=63 got=%d", relBitIdx)
	}

	// first bit 2nd chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(64)
	if relIntIdx != 1 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 0 {
		t.Errorf("expected=0 got=%d", relBitIdx)
	}

	// 2nd bit 2nd chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(65)
	if relIntIdx != 1 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 1 {
		t.Errorf("expected=1 got=%d", relBitIdx)
	}

	// last bit 2nd chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(64 + 63)
	if relIntIdx != 1 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 63 {
		t.Errorf("expected=63 got=%d", relBitIdx)
	}

	// first bit last chunk
	relIntIdx, relBitIdx, _ = storage.getIdx(64 + 64)
	if relIntIdx != 2 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 0 {
		t.Errorf("expected=0 got=%d", relBitIdx)
	}

	// very last available entry
	relIntIdx, relBitIdx, _ = storage.getIdx(64 + 64 + 63)
	if relIntIdx != 2 {
		t.Errorf("expected=1 got=%d", relIntIdx)
	}
	if relBitIdx != 63 {
		t.Errorf("expected=63 got=%d", relBitIdx)
	}
}

func Test_getConvertedInt(t *testing.T) {
	expectedZero := "0000000000000000000000000000000000000000000000000000000000000000"
	if getConvertedInt(0) != expectedZero {
		t.Errorf("expected %s to equal %s", getConvertedInt(0), expectedZero)
	}

	expectedOne := "0000000000000000000000000000000000000000000000000000000000000001"
	if getConvertedInt(1) != expectedOne {
		t.Errorf("expected %s to equal %s", getConvertedInt(1), expectedOne)
	}

	expectedTwo := "0000000000000000000000000000000000000000000000000000000000000010"
	if getConvertedInt(2) != expectedTwo {
		t.Errorf("expected %s to equal %s", getConvertedInt(2), expectedTwo)
	}
}

func TestStorage_AddItem(t *testing.T) {

	input, storage := storageBeforeEach()

	// first possible item
	hash := uint64(0)
	err := storage.AddItem(hash)
	if err != nil {
		t.Error(err)
	}
	if !storage.PotentiallyKnowItem(hash) {
		t.Errorf("failed to find added item")
	}

	// 2nd possible item
	hash = uint64(1)
	err = storage.AddItem(hash)
	if err != nil {
		t.Error(err)
	}
	if !storage.PotentiallyKnowItem(hash) {
		t.Errorf("failed to find added item")
	}

	// last required item
	hash = input.MinEntryCount
	err = storage.AddItem(hash)
	if err != nil {
		t.Error(err)
	}

	if !storage.PotentiallyKnowItem(hash) {
		t.Errorf("failed to find added item")
	}
}
