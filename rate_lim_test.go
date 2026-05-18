package rateLimiter
import "testing"
func TestLimiter(t *testing.T) {
	lim := NewRateLimiter(2, 100, 1000, 10000, 100000, 100000, 1000000, 100)

	ip := "1.1.1.1"

	if !lim.Trust(ip) {
		t.Error("should allow")
	}
	if !lim.Trust(ip) {
		t.Error("should allow")
	}
	if lim.Trust(ip) {
		t.Error("should block by rps")
	}
}
func TestLimiterDifferentIps(test *testing.T) {
	limit := NewRateLimiter(1, 100, 1000, 10000, 100000, 100000, 1000000, 100)
	if !limit.Adobe("1.1.1.1") {
		test.Error("should")
	}
	if !limit.Adobe("2.2.2.2") {
		test.Error("different ips")
	}
}
func TestUploadLimit(test *testing.T) {
	lim := NewRateLimiter(
		100,
		100,
		100,
		100,
		10,
		100,
		1000,
		100,
	)
	ok := lim.CheckUpload("1.1.1.1", 100)
	if ok {
		test.Error("must exceed upload (TestUploadLimit)")
	}
}

func TestDownloadLimit(test *testing.T) {
	lim := NewRateLimiter(
		100,
		100,
		100,
		100,
		1000,
		10,
		1000,
		100,
	)
	ok := lim.CheckDownload("1.1.1.1", 100)
	if ok {
		test.Error("must exceed download (TestDownloadLimit)")
	}
}
func TestMaxConnections(test *testing.T) {
	lim := NewRateLimiter(
		10,
		100,
		1000,
		10000,
		100000,
		100000,
		1000000,
		1,
	)
	if lim.CheckMaxConnections(1) {
		test.Error("must fail (TestMaxConnections)")
	}
}
func TestSubnetLimiter(test *testing.T) {
	lim := NewRateLimiter(1, 100, 1000, 10000, 1000, 1000, 1000, 100)
	if !lim.Trust("192.168.1.1") {
		test.Error("first should pass (TestSubnetLimiter)")
	}
	if !lim.Trust("192.168.1.2") {
		test.Error("second subnet should pass (TestSubnetLimiter)")
	}
}
func TestSetLimits(test *testing.T) {
	lim := NewRateLimiter(1, 1, 1, 1, 1, 1, 1, 1)
	lim.SetLimits(10, 20, 30, 40)
	if lim.rps != 10 {
		test.Error("set limit fail (TestSetLimits)")
	}
} 