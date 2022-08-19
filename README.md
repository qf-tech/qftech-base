# qftech-base
some base lib for go

## features
### log: 
the logger that use context inject traceID  
rotate the log file  
log format support json and commone text line  

### crypt(encrypt & decrypt)
use AES256-CBC and PBKDF2 to encrypt or decrypt data. For example, encrypt and decrypt sensitive configuration items  

### sign(generate signature and verify)
use RSA to generate signature of the file hash(sha256), also you can verify the signature

*Notice*: you can find their usages in example directory
