// Results behind the benchmark chart and tables. Rerun with the commands on the Benchmarks page and update together.

export const environment = 'Apple M4 Max, Go 1.27.0, median of 5 runs';

export type Adapter = 'nethttp' | 'fasthttp';

export interface Result {
	ns: number;
	allocs: number;
}

export interface Comparison {
	scenario: string;
	v4: Result;
	v5: Result;
}

/** v4 against v5 on the same harness: no-op handlers and middleware, one request per op. */
export const comparison: Record<Adapter, Comparison[]> = {
	nethttp: [
		{ scenario: 'Static route', v4: { ns: 20.0, allocs: 0 }, v5: { ns: 19.9, allocs: 0 } },
		{ scenario: 'Static, 5 segments', v4: { ns: 28.5, allocs: 0 }, v5: { ns: 30.2, allocs: 0 } },
		{ scenario: 'One param', v4: { ns: 87.0, allocs: 4 }, v5: { ns: 93.6, allocs: 2 } },
		{ scenario: 'Three params + regexp', v4: { ns: 150.1, allocs: 4 }, v5: { ns: 148.5, allocs: 2 } },
		{ scenario: 'Static + 3 middleware', v4: { ns: 170.4, allocs: 6 }, v5: { ns: 23.8, allocs: 0 } },
		{ scenario: 'One param + 3 middleware', v4: { ns: 238.4, allocs: 10 }, v5: { ns: 100.8, allocs: 2 } },
		{ scenario: '404 Not Found', v4: { ns: 20.0, allocs: 0 }, v5: { ns: 20.3, allocs: 0 } },
		{ scenario: '405 Method Not Allowed', v4: { ns: 60.8, allocs: 3 }, v5: { ns: 55.8, allocs: 2 } },
	],
	fasthttp: [
		{ scenario: 'Static route', v4: { ns: 30.0, allocs: 1 }, v5: { ns: 16.9, allocs: 0 } },
		{ scenario: 'Static, 5 segments', v4: { ns: 39.0, allocs: 1 }, v5: { ns: 26.2, allocs: 0 } },
		{ scenario: 'One param', v4: { ns: 50.1, allocs: 3 }, v5: { ns: 38.5, allocs: 0 } },
		{ scenario: 'Three params + regexp', v4: { ns: 110.2, allocs: 3 }, v5: { ns: 87.4, allocs: 0 } },
		{ scenario: 'Static + 3 middleware', v4: { ns: 173.1, allocs: 7 }, v5: { ns: 18.3, allocs: 0 } },
		{ scenario: 'One param + 3 middleware', v4: { ns: 192.8, allocs: 9 }, v5: { ns: 39.5, allocs: 0 } },
		{ scenario: '404 Not Found', v4: { ns: 29.7, allocs: 1 }, v5: { ns: 15.3, allocs: 0 } },
		{ scenario: '405 Method Not Allowed', v4: { ns: 89.0, allocs: 3 }, v5: { ns: 69.5, allocs: 1 } },
	],
};

export interface SuiteResult {
	benchmark: string;
	what: string;
	ns: number;
	bytes: number;
	allocs: number;
}

/** v5's own benchmark suite, `go test -bench . ./...`. */
export const suite: Record<Adapter, SuiteResult[]> = {
	nethttp: [
		{ benchmark: 'Static/1', what: 'One literal segment', ns: 15.0, bytes: 0, allocs: 0 },
		{ benchmark: 'Static/5', what: 'Five literal segments', ns: 31.6, bytes: 0, allocs: 0 },
		{ benchmark: 'Static/20', what: 'Twenty literal segments', ns: 98.8, bytes: 0, allocs: 0 },
		{ benchmark: 'Param/1', what: 'One param', ns: 90.3, bytes: 336, allocs: 2 },
		{ benchmark: 'Param/5', what: 'Five params', ns: 132.2, bytes: 336, allocs: 2 },
		{ benchmark: 'Param/20', what: 'Twenty params', ns: 883.4, bytes: 2137, allocs: 7 },
		{ benchmark: 'Regexp/1', what: 'One regexp param', ns: 157.6, bytes: 336, allocs: 2 },
		{ benchmark: 'Regexp/5', what: 'Five regexp params', ns: 434.8, bytes: 336, allocs: 2 },
		{ benchmark: 'Middleware', what: 'Literal route behind 3 middleware', ns: 23.1, bytes: 0, allocs: 0 },
		{ benchmark: 'NotFound', what: '404', ns: 18.9, bytes: 0, allocs: 0 },
		{ benchmark: 'MethodNotAllowed', what: '405 with an Allow header', ns: 55.7, bytes: 32, allocs: 2 },
		{ benchmark: 'Mount', what: 'Request through a mounted router', ns: 83.0, bytes: 448, allocs: 1 },
		{ benchmark: 'API', what: 'Average over a 70-route API', ns: 95.5, bytes: 226, allocs: 1.3 },
	],
	fasthttp: [
		{ benchmark: 'Static/1', what: 'One literal segment', ns: 11.6, bytes: 0, allocs: 0 },
		{ benchmark: 'Static/5', what: 'Five literal segments', ns: 27.1, bytes: 0, allocs: 0 },
		{ benchmark: 'Static/20', what: 'Twenty literal segments', ns: 93.3, bytes: 0, allocs: 0 },
		{ benchmark: 'Param/1', what: 'One param', ns: 33.0, bytes: 0, allocs: 0 },
		{ benchmark: 'Param/5', what: 'Five params', ns: 66.6, bytes: 0, allocs: 0 },
		{ benchmark: 'Param/20', what: 'Twenty params', ns: 250.0, bytes: 0, allocs: 0 },
		{ benchmark: 'Regexp/1', what: 'One regexp param', ns: 89.7, bytes: 0, allocs: 0 },
		{ benchmark: 'Regexp/5', what: 'Five regexp params', ns: 353.1, bytes: 0, allocs: 0 },
		{ benchmark: 'Middleware', what: 'Literal route behind 3 middleware', ns: 17.1, bytes: 0, allocs: 0 },
		{ benchmark: 'NotFound', what: '404', ns: 14.3, bytes: 0, allocs: 0 },
		{ benchmark: 'MethodNotAllowed', what: '405 with an Allow header', ns: 65.9, bytes: 16, allocs: 1 },
		{ benchmark: 'Mount', what: 'Request through a mounted router, URI parsing included', ns: 130.2, bytes: 0, allocs: 0 },
		{ benchmark: 'API', what: 'Average over a 70-route API', ns: 53.1, bytes: 0, allocs: 0 },
	],
};
