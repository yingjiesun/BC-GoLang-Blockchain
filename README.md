YS 4/16/2018

Updates: 
- Separated files to main, block, blockchain, transaction, tcp, http;
- Updated block to use []Transaction as data
- Updated calculateHash(), generateBlock() to use []Transaction
- Update tcp to take input and convert to []Transaction

HTTP Server port: 8881
TCP Server port: 9991

TODO: http taking jason obj and convert to []Transaction
