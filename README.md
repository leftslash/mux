#### Note
This is a work in progress.  Use at your own risk.

#### Fixes
- Normalize regex names in parser (e.g. ParserDirRegex)
- Fix json funcs to write body after errors detected, not vice versa

#### Further Reading
- [Good HTTP Citizen][1]
- [Exposing Go on the Internet][2]

#### Enhancements
- Make Use() variadic and return \*Router
- Add appropriate secure cookie options to auth
- Add signal handler to shutdwon gracefully
- Add static utility method
- Add flags, usage, version to serve-static
- Add comments
- Add documentation

#### Other Stuff
- CSRF?
- CORS?
- Read more at [OWASP][3]

[1]: https://blog.devgenius.io/writing-custom-http-servers-as-good-goolang-citizens-ea6b2ebfc05f
[2]: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
[3]: https://cheatsheetseries.owasp.org/
