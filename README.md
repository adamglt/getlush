# `getlush`

Batch downloader for payslips from Hilan.

### Usage

1. Steal your own browser cookie:
    - Open your browser developer console and switch to the `Network` tab.
    - Login to Hilan as you normally would and view your latest payslip.
    - Find a request that looks like `PaySlip2022-01.pdf?Date=01/01/2022&UserId=12341231231`.
    - Go to `Request Headers` on the bottom right and copy the `Cookie` header.
    - Save the value to a file, e.g. `hilan.cookie` under `bin` directory.

2. Create executable file by running `make`

3. Run `getlush` with your details - for example:  
    `$ ./bin/getlush -emp 209935777 -from 2015-01 -to 2022-01`


### Config

```
$ ./getlush -help
Usage of ./getlush:
  -cookie string
        Path to cookie file (default "hilan.cookie")
  -emp string
        Employee citizenship ID number (not userid) [required]
  -from value
        First payslip to fetch (YYYY-MM) [required]
  -org string
        Parent organization ID (default "9133")
  -out string
        Directory path for fetched pdfs (default "getlush_out/")
  -t duration
        Single request timeout (default 10s)
  -to value
        Last payslip to fetch (YYYY-MM) [required]
  -url string
        Hilan's base URL (default "https://traiana.net.hilan.co.il/")

```
