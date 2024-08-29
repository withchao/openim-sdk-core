package checker

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	TotalCount   int
	CorrectCount int
	IsEqual      bool
}

func NewCounter(total, correct int, isEqual bool) *Counter {
	return &Counter{
		TotalCount:   total,
		CorrectCount: correct,
		IsEqual:      isEqual,
	}
}

type CounterChecker[T any, K comparable] struct {
	// CheckName must start with 'check' and be named with a small camel hump,
	// followed by the name of the indicator that needs to be checked,
	// and it will be assigned to checkNumName.
	// e.g. checkGroupNum: checkNumName=GroupNum
	CheckName       string
	checkNumName    string // used for printing logs
	CheckerKeyName  string // used for printing logs
	GoroutineLimit  int
	GetTotalCount   func(ctx context.Context, t T) (int, error) // get now total count
	CalCorrectCount func(key K) int                             // return correct num
	LoopSlice       []T                                         // circular slicing
	GetKey          func(t T) K                                 // get checkers key from a type
}

func (c *CounterChecker[T, K]) Init() {
	c.CheckName = stringutil.LowerFirst(c.CheckName)
	c.checkNumName = strings.TrimPrefix(c.CheckName, "check")
}

func (c *CounterChecker[T, K]) Check(ctx context.Context) error {
	defer decorator.FuncLogSkip(ctx, 1)()

	c.Init()

	var (
		gr, cctx = reerrgroup.WithContext(ctx, c.GoroutineLimit)
		total    atomic.Int64
		progress atomic.Int64

		checkers = make(map[K]*Counter, len(sdk.TestSDKs))
		mapLock  = sync.RWMutex{}
	)

	total.Add(int64(len(c.LoopSlice)))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)

	for _, t := range c.LoopSlice {
		t := t
		gr.Go(func() error {
			key := c.GetKey(t)
			correctNum := c.CalCorrectCount(key)
			totalNum, err := c.GetTotalCount(ctx, t)
			if err != nil {
				return err
			}
			isEqual := totalNum == correctNum
			if !isEqual {
				mapLock.Lock()
				checkers[key] = NewCounter(totalNum, correctNum, isEqual)
				mapLock.Unlock()
			}
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}

	if len(checkers) != 0 {
		err := errs.New(fmt.Sprintf("%s un correct!", stringutil.CamelCaseToSpaceSeparated(c.CheckName))).Wrap()
		for k, ck := range checkers {
			log.ZWarn(ctx, fmt.Sprintf("%s un correct", stringutil.CamelCaseToSpaceSeparated(c.checkNumName)),
				err, c.CheckerKeyName, k, c.checkNumName, ck.TotalCount, "correct num", ck.CorrectCount)
		}
		InsertToErrChan(ctx, err)
	} else {
		log.ZInfo(ctx, fmt.Sprintf("%s success", stringutil.CamelCaseToSpaceSeparated(c.CheckName)))
	}
	return nil
}

func (c *CounterChecker[T, K]) LoopCheck(ctx context.Context) error {
	defer decorator.FuncLogSkip(ctx, 1)()

	c.Init()

	var (
		gr, cctx = reerrgroup.WithContext(ctx, c.GoroutineLimit)
		total    atomic.Int64
		progress atomic.Int64

		checkers = make(map[K]*Counter, len(sdk.TestSDKs))
	)

	total.Add(int64(len(c.LoopSlice)))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)

	for _, t := range c.LoopSlice {
		t := t
		checkCount := 0
		isEqual := false
		gr.Go(func() error {
			key := c.GetKey(t)
			correctNum := c.CalCorrectCount(key)
			for !isEqual {

				totalNum, err := c.GetTotalCount(ctx, t)
				if err != nil {
					return err
				}
				isEqual = totalNum == correctNum
				if !isEqual {
					checkCount++

					logStr := fmt.Sprintf("check num:%d, %s un correct. %s: %s,%s: %d, %s: %d",
						checkCount, stringutil.CamelCaseToSpaceSeparated(c.checkNumName),
						c.CheckerKeyName, key, c.checkNumName, totalNum, "correct num", correctNum)

					fmt.Println(logStr)
					log.ZWarn(ctx, fmt.Sprintf("check num:%d, %s un correct",
						checkCount, stringutil.CamelCaseToSpaceSeparated(c.checkNumName)),
						nil, c.CheckerKeyName, key, c.checkNumName, totalNum, "correct num", correctNum)
				}
				time.Sleep(config.CheckWaitSec * time.Second)
			}
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}

	if len(checkers) != 0 {
		err := errs.New(fmt.Sprintf("%s un correct!", stringutil.CamelCaseToSpaceSeparated(c.CheckName))).Wrap()
		for k, ck := range checkers {
			log.ZWarn(ctx, fmt.Sprintf("%s un correct", stringutil.CamelCaseToSpaceSeparated(c.checkNumName)),
				err, c.CheckerKeyName, k, c.checkNumName, ck.TotalCount, "correct num", ck.CorrectCount)
		}
		InsertToErrChan(ctx, err)
	} else {
		log.ZInfo(ctx, fmt.Sprintf("%s success", stringutil.CamelCaseToSpaceSeparated(c.CheckName)))
	}
	return nil
}