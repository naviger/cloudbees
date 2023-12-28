package main

import (
	"time"

	"google.golang.org/grpc/metadata"
)

func GetMetadataKey(md metadata.MD, key string, deflt any) any {
	arr := md.Get(key)
	if len(arr) > 0 {
		return arr[0]
	} else {
		return deflt
	}
}

func GetFirstOpenSeat(dt time.Time) SeatS {
	txn := db.Txn(false)
	defer txn.Abort()

	found := false
	var st SeatS = SeatS{}

	for !found {
		it, _ := txn.LowerBound("Seat", "trainId", dt.Format("20060102"))
		for obj := it.Next(); obj != nil; obj = it.Next() {
			st = obj.(SeatS)
			if st.Status == "vacant" {
				found = true
				return st
			}
		}

	}
	return SeatS{}
}
