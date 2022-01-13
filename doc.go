//
// 	/a        ok, matches /a but not /a/
// 	/a/*      ok, matches /a/ but not /a
// 	/a/       ok, matches /a and /a/
//
// 	a         error, leading slash required
// 	/         ok, matches anything
// 	/*        ok, matches anything
// 	/a*       error, wildcard must follow slash
// 	/a/*      ok, matches /a/ but not /a
// 	/*a       error, wildcard must be last
// 	/*/       error, wildcard must be last
// 	/a*/      error, wildcard must be last
//
// 	/ a       error, patterns must not contain spaces
// 	/ a/      error, patterns must not contain spaces
// 	/a        error, patterns must not contain spaces
// 	/a /      error, patterns must not contain spaces
// 	/{ }      error, patterns must not contain spaces
// 	/{id }    error, patterns must not contain spaces
// 	/{ id}    error, patterns must not contain spaces
// 	/{i d}    error, patterns must not contain spaces
//
//
// 	{id}      error, leading slash required
// 	/{id}     ok, matches /5 but not /5/
// 	/{id}*    error, wildcard must follow slash
// 	/{id}/    ok, matches /5 and /5/, with id = 5 for both
// 	/{id}/*   ok, matches /5/ and /5/5 but not /5
//
// 	/a{id}    ok
// 	/{id}a    ok
// 	/a{id}a   ok
//
// 	/{}       error, empty parameter name
// 	/{{}}     error, nested parameters
// 	/{{id}}   error, nested parameters
// 	/{a{id}   error, nested parameters
// 	/{id}a}   error, nested parameters
// 	/{a{id}b} error, nested parameters
//
// 	/{/}      error, parameter cannot contain slash
// 	/{/id}    error, cannot contain slash
// 	/{id/}    error, cannot contain slash
// 	/{i/d}    error, cannot contain slash
//
// 	//        ok, same as /
// 	//a       ok, same as /a
// 	/a//      ok, same as /a/
// 	//a//     ok, same as /a/
//
package mux
