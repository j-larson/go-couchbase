package couchbase

import "testing"

func TestSingleNode(t *testing.T) {
	jsonInput := `{"rev":66,"nodesExt":[{"services":{"mgmt":8091,"mgmtSSL":18091,"indexAdmin":9100,"indexScan":9101,"indexHttp":9102,"indexStreamInit":9103,"indexStreamCatchup":9104,"indexStreamMaint":9105,"indexHttps":19102,"kv":11210,"kvSSL":11207,"capi":8092,"capiSSL":18092,"projector":9999,"n1ql":8093,"n1qlSSL":18093},"thisNode":true}],"clusterCapabilitiesVer":[1,0],"clusterCapabilities":{"n1ql":["enhancedPreparedStatements"]}}`

	poolServices, err := ParsePoolServices(jsonInput)
	if err != nil {
		t.Fatalf("Unable to parse json: %v", err)
	}
	if poolServices == nil {
		t.Fatalf("Parse produced no result")
	}
	if len(poolServices.NodesExt) != 1 {
		t.Fatalf("Expected nodesExt of length 1, got %d", len(poolServices.NodesExt))
	}
	if poolServices.NodesExt[0].Hostname != "" {
		t.Fatalf("Expected empty hostname, got %s", poolServices.NodesExt[0].Hostname)
	}
	if poolServices.NodesExt[0].Services["kv"] != 11210 {
		t.Fatalf("Expected kv port 11210, got %d", poolServices.NodesExt[0].Services["kv"])
	}
	if poolServices.NodesExt[0].Services["kvSSL"] != 11207 {
		t.Fatalf("Expected kvSSL port 11207, got %d", poolServices.NodesExt[0].Services["kvSSL"])
	}

	// Should succeed.
	target := "127.0.0.1:11210"
	res, err := MapKVtoSSL(target, poolServices)
	if err != nil {
		t.Fatalf("Mapping target %s, expected success, got error: %v", target, err)
	}
	expected := "127.0.0.1:11207"
	if res != expected {
		t.Fatalf("Mapping target %s, expected %s, got %s", target, expected, res)
	}

	// No port.
	target = "127.0.0.1"
	res, err = MapKVtoSSL(target, poolServices)
	if err == nil {
		t.Fatalf("Mapping target %s, expected failure, got success: %s", target, res)
	}

	// Bad KV port.
	target = "127.0.0.1:11111"
	res, err = MapKVtoSSL(target, poolServices)
	if err == nil {
		t.Fatalf("Mapping target %s, expected failure, got success: %s", target, res)
	}
}

func TestMultiNode(t *testing.T) {
	jsonInput := `{"rev":32,"nodesExt":[{"services":{"mgmt":8091,"mgmtSSL":18091,"fts":8094,"ftsSSL":18094,"indexAdmin":9100,"indexScan":9101,"indexHttp":9102,"indexStreamInit":9103,"indexStreamCatchup":9104,"indexStreamMaint":9105,"indexHttps":19102,"capiSSL":18092,"capi":8092,"kvSSL":11299,"projector":9999,"kv":11298,"moxi":11211},"hostname":"172.23.123.101"},{"services":{"mgmt":8091,"mgmtSSL":18091,"indexAdmin":9100,"indexScan":9101,"indexHttp":9102,"indexStreamInit":9103,"indexStreamCatchup":9104,"indexStreamMaint":9105,"indexHttps":19102,"capiSSL":18092,"capi":8092,"kvSSL":11207,"projector":9999,"kv":11210,"moxi":11211,"n1ql":8093,"n1qlSSL":18093},"thisNode":true,"hostname":"172.23.123.102"}]}`

	poolServices, err := ParsePoolServices(jsonInput)
	if err != nil {
		t.Fatalf("Unable to parse json: %v", err)
	}
	if poolServices == nil {
		t.Fatalf("Parse produced no result")
	}
	if len(poolServices.NodesExt) != 2 {
		t.Fatalf("Expected nodesExt of length 2, got %d", len(poolServices.NodesExt))
	}
	if poolServices.NodesExt[0].Services["kv"] != 11298 {
		t.Fatalf("Expected kv port 11298, got %d", poolServices.NodesExt[0].Services["kv"])
	}
	if poolServices.NodesExt[1].Services["kvSSL"] != 11207 {
		t.Fatalf("Expected kvSSL port 11207, got %d", poolServices.NodesExt[1].Services["kvSSL"])
	}

	// Should succeed.
	target := "172.23.123.102:11210"
	res, err := MapKVtoSSL(target, poolServices)
	if err != nil {
		t.Fatalf("Mapping target %s, expected success, got error: %v", target, err)
	}
	expected := "172.23.123.102:11207"
	if res != expected {
		t.Fatalf("Mapping target %s, expected %s, got %s", target, expected, res)
	}

	// No such host.
	target = "172.23.123.999:11210"
	res, err = MapKVtoSSL(target, poolServices)
	if err == nil {
		t.Fatalf("Mapping target %s, expected failure, got success: %s", target, res)
	}

	// Bad KV port.
	target = "172.23.123.101:11111"
	res, err = MapKVtoSSL(target, poolServices)
	if err == nil {
		t.Fatalf("Mapping target %s, expected failure, got success: %s", target, res)
	}
}

func TestIPv6Node(t *testing.T) {
	jsonInput := `{"rev":32,"nodesExt":[{"services":{"mgmt":8091,"mgmtSSL":18091,"fts":8094,"ftsSSL":18094,"indexAdmin":9100,"indexScan":9101,"indexHttp":9102,"indexStreamInit":9103,"indexStreamCatchup":9104,"indexStreamMaint":9105,"indexHttps":19102,"capiSSL":18092,"capi":8092,"kvSSL":11299,"projector":9999,"kv":11298,"moxi":11211},"hostname":"[DEAD::BEEF]"},{"services":{"mgmt":8091,"mgmtSSL":18091,"indexAdmin":9100,"indexScan":9101,"indexHttp":9102,"indexStreamInit":9103,"indexStreamCatchup":9104,"indexStreamMaint":9105,"indexHttps":19102,"capiSSL":18092,"capi":8092,"kvSSL":11207,"projector":9999,"kv":11210,"moxi":11211,"n1ql":8093,"n1qlSSL":18093},"thisNode":true,"hostname":"[FEED::DEED]"}]}`

	poolServices, err := ParsePoolServices(jsonInput)
	if err != nil {
		t.Fatalf("Unable to parse json: %v", err)
	}
	if poolServices == nil {
		t.Fatalf("Parse produced no result")
	}
	if len(poolServices.NodesExt) != 2 {
		t.Fatalf("Expected nodesExt of length 2, got %d", len(poolServices.NodesExt))
	}
	if poolServices.NodesExt[0].Services["kv"] != 11298 {
		t.Fatalf("Expected kv port 11298, got %d", poolServices.NodesExt[0].Services["kv"])
	}
	if poolServices.NodesExt[1].Services["kvSSL"] != 11207 {
		t.Fatalf("Expected kvSSL port 11207, got %d", poolServices.NodesExt[1].Services["kvSSL"])
	}

	// Should succeed.
	target := "[FEED::DEED]:11210"
	res, err := MapKVtoSSL(target, poolServices)
	if err != nil {
		t.Fatalf("Mapping target %s, expected success, got error: %v", target, err)
	}
	expected := "[FEED::DEED]:11207"
	if res != expected {
		t.Fatalf("Mapping target %s, expected %s, got %s", target, expected, res)
	}

	// Bad KV port.
	target = "[DEAD::BEEF]:11111"
	res, err = MapKVtoSSL(target, poolServices)
	if err == nil {
		t.Fatalf("Mapping target %s, expected failure, got success: %s", target, res)
	}
}
