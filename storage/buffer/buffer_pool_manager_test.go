// this code is from https://github.com/brunocalza/go-bustub
// there is license and copyright notice in licenses/go-bustub dir

package buffer

import (
	"crypto/rand"
	testingpkg "github.com/ryogrid/bltree-go-for-embedding/testing/testing_assert"
	"testing"

	"github.com/ryogrid/bltree-go-for-embedding/common"
	"github.com/ryogrid/bltree-go-for-embedding/storage/disk"
	"github.com/ryogrid/bltree-go-for-embedding/types"
)

func TestBinaryData(t *testing.T) {
	common.TempSuppressOnMemStorageMutex.Lock()
	common.TempSuppressOnMemStorage = true

	poolSize := uint32(10)

	dm := disk.NewDiskManagerTest()
	bpm := NewBufferPoolManager(poolSize, dm)

	page0 := bpm.NewPage()

	// Scenario: The buffer pool is empty. We should be able to create a new page.
	testingpkg.Equals(t, types.PageID(0), page0.GetPageId())

	// Generate random binary data
	randomBinaryData := make([]byte, common.PageSize)
	rand.Read(randomBinaryData)

	// Insert terminal characters both in the middle and at end
	randomBinaryData[common.PageSize/2] = '0'
	randomBinaryData[common.PageSize-1] = '0'

	var fixedRandomBinaryData [common.PageSize]byte
	copy(fixedRandomBinaryData[:], randomBinaryData[:common.PageSize])

	// Scenario: Once we have a page, we should be able to read and write content.
	page0.Copy(0, randomBinaryData)
	testingpkg.Equals(t, fixedRandomBinaryData, *page0.Data())

	// Scenario: We should be able to create new pages until we fill up the buffer pool.
	for i := uint32(1); i < poolSize; i++ {
		p := bpm.NewPage()
		testingpkg.Equals(t, types.PageID(i), p.GetPageId())
	}

	// Scenario: After unpinning pages {0, 1, 2, 3, 4} and pinning another 4 new pages,
	// there would still be one cache frame left for reading page 0.
	for i := 0; i < 5; i++ {
		testingpkg.Ok(t, bpm.UnpinPage(types.PageID(i), true))
		bpm.FlushPage(types.PageID(i))
	}
	for i := 0; i < 4; i++ {
		p := bpm.NewPage()
		bpm.UnpinPage(p.GetPageId(), false)
	}

	// Scenario: We should be able to fetch the data we wrote a while ago.
	page0 = bpm.FetchPage(types.PageID(0))
	testingpkg.Equals(t, fixedRandomBinaryData, *page0.Data())
	testingpkg.Ok(t, bpm.UnpinPage(types.PageID(0), true))

	common.TempSuppressOnMemStorage = false
	dm.ShutDown()
	common.TempSuppressOnMemStorageMutex.Unlock()
}

func TestSample(t *testing.T) {
	common.TempSuppressOnMemStorageMutex.Lock()
	common.TempSuppressOnMemStorage = true

	poolSize := uint32(10)

	dm := disk.NewDiskManagerTest()
	bpm := NewBufferPoolManager(poolSize, dm)

	page0 := bpm.NewPage()

	// Scenario: The buffer pool is empty. We should be able to create a new page.
	testingpkg.Equals(t, types.PageID(0), page0.GetPageId())

	// Scenario: Once we have a page, we should be able to read and write content.
	page0.Copy(0, []byte("Hello"))
	testingpkg.Equals(t, [common.PageSize]byte{'H', 'e', 'l', 'l', 'o'}, *page0.Data())

	// Scenario: We should be able to create new pages until we fill up the buffer pool.
	for i := uint32(1); i < poolSize; i++ {
		p := bpm.NewPage()
		testingpkg.Equals(t, types.PageID(i), p.GetPageId())
	}

	// Scenario: After unpinning pages {0, 1, 2, 3, 4} and pinning another 4 new pages,
	// there would still be one cache frame left for reading page 0.
	for i := 0; i < 5; i++ {
		testingpkg.Ok(t, bpm.UnpinPage(types.PageID(i), true))
		bpm.FlushPage(types.PageID(i))
	}
	for i := 0; i < 4; i++ {
		bpm.NewPage()
	}
	// Scenario: We should be able to fetch the data we wrote a while ago.
	page0 = bpm.FetchPage(types.PageID(0))
	testingpkg.Equals(t, [common.PageSize]byte{'H', 'e', 'l', 'l', 'o'}, *page0.Data())

	// Scenario: If we unpin page 0 and then make a new page, all the buffer pages should
	// now be pinned. Fetching page 0 should fail.
	testingpkg.Ok(t, bpm.UnpinPage(types.PageID(0), true))

	p := bpm.NewPage()
	testingpkg.Equals(t, types.PageID(14), p.GetPageId())

	common.TempSuppressOnMemStorage = false
	dm.ShutDown()
	common.TempSuppressOnMemStorageMutex.Unlock()
}
