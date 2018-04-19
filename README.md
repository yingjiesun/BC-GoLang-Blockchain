YS 4/19/2018<BR>
Added BLOCKINTERVAL=5 to .env<BR>

YS 4/18/2018<BR>
Added temp_trans to main.go, to hold temp transactions that are not saved in blockchain yet<BR>

YS 4/16/2018<BR>
Break project to multiple files<BR>
Add transaction.go<BR>
Modify block.go to have transactions array<BR>
Modify TCP to accept new blocks in new format<BR>
Add nounce to block<BR>
Add DIFFICULTY in .env, 5 leading 0s now<BR>
Add POW<BR>

YS 4/14/2018<BR>
Updates:<BR> 
Separated files to main, block, http, transaction;<BR>
Both Http TCP server starts when app is launched<BR>
HTTP Server port: 8881<BR>
TCP Server port: 9991<BR>