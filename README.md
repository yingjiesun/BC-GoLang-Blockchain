YS 4/16/2018
Break project to multiple files
Add transaction.go
Modify block.go to have transactions array
Modify TCP to accept new blocks in new format
Add nounce to block
Add DIFFICULTY in .env, 5 leading 0s now
Add POW

YS 4/14/2018
Updates: 
Separated files to main, block, http, transaction;
Both Http TCP server starts when app is launched
HTTP Server port: 8881
TCP Server port: 9991
