package setplatformname

import (
	"context"
	"encoding/csv"
	"io"
	"leishen/pkg/edb"
	"leishen/pkg/etypes"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

const (
	csv_column_count = 2
	parallelism      = 16
)

type nameSet map[string]struct{}

func setPlatformName(contractDB, fltxDB ethdb.Database, csvFilePath string) error {
	platformName, err := getPlatformNameFromCSV(csvFilePath)
	if err != nil {
		return err
	}

	eoaMap := make(map[common.Address]struct{})
	accountNameSet := make(map[common.Address]nameSet)
	bar := progressbar.Default(int64(len(platformName)), "Build the contract creation diagram")
	for address, name := range platformName {
		eoaPath := getEOAPath(contractDB, address)

		eoaMap[eoaPath[len(eoaPath)-1]] = struct{}{}

		for _, addr := range eoaPath {
			if c, ok := accountNameSet[addr]; !ok {
				accountNameSet[addr] = map[string]struct{}{name: {}}
			} else {
				c[name] = struct{}{}
			}
		}
		bar.Add(1)
	}

	bar = progressbar.Default(int64(len(eoaMap)), "Set platform name")

	grounp, ctx := errgroup.WithContext(context.Background())
	eoaCh := make(chan common.Address, parallelism)
	for i := 0; i < parallelism; i++ {
		grounp.Go(func() error {
			return updataSubcontractByChan(ctx, contractDB, fltxDB, accountNameSet, eoaCh)
		})
	}

	grounp.Go(func() error {
		for eoa := range eoaMap {
			select {
			case <-ctx.Done():
				return nil

			case eoaCh <- eoa:
				bar.Add(1)
			}
		}
		close(eoaCh)
		return nil
	})

	return grounp.Wait()
}

func getPlatformNameFromCSV(fileName string) (map[common.Address]string, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0400)
	if err != nil {
		return nil, errCSVOpen(err)
	}
	defer file.Close()

	platformName := make(map[common.Address]string)
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			return platformName, nil
		}

		if err != nil {
			return nil, errCSVRead(fileName, err)
		}

		trimmed_record := trimSliceSpace(record)
		if len(trimmed_record) != csv_column_count {
			return nil, errCSVItemNumber(csv_column_count, len(trimmed_record))
		}

		address := common.HexToAddress(record[0])
		name := record[1]
		if name == "" {
			continue
		}

		platformName[address] = name
	}
}

func getEOAPath(contractDB ethdb.Database, addr common.Address) (eoaPath []common.Address) {
	for {
		eoaPath = append(eoaPath, addr)

		contract := edb.ReadContract(contractDB, &addr)
		if contract == nil {
			return
		}

		addr = contract.Creator
	}
}

func updataSubcontractByChan(ctx context.Context, contractDB, fltxDB ethdb.Database, accountNameSet map[common.Address]nameSet, addr chan common.Address) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case eoa_addr, ok := <-addr:
			if !ok {
				return nil
			}

			err := updataSubcontract(contractDB, fltxDB, accountNameSet, eoa_addr, accountNameSet[eoa_addr], true)
			if err != nil {
				return err
			}
		}
	}
}

func updataSubcontract(contractDB, fltxDB ethdb.Database, accountNameSet map[common.Address]nameSet, addr common.Address, name map[string]struct{}, isEOA bool) error {
	platformName := &etypes.PlatformName{
		Address: addr,
		IsEOA:   isEOA,
		NameMap: name,
	}

	edb.WritePlatformName(fltxDB, platformName)

	iter := edb.GetSubContractIterator(contractDB, addr)
	defer iter.Release()
	for {
		childAddr, ok := iter.Next()
		if !ok {
			return nil
		}

		if account, ok := accountNameSet[*childAddr]; ok {
			updataSubcontract(contractDB, fltxDB, accountNameSet, *childAddr, account, false)
			continue
		}

		updataSubcontract(contractDB, fltxDB, accountNameSet, *childAddr, name, false)
	}
}
