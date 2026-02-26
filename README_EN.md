# KobackupCipherTool-go

A Go implementation of Kobackup backup file decryption tool. Used to parse and decrypt encrypted backup files from some phones.

## Features

- **checkhash**: Verify backup file integrity
  - Uses backup password and HMAC-SHA256 algorithm to verify signatures
  - Supports parsing checkMsgV3 fields from `info.xml`

- **decrypt**: Decrypt backup files
  - Uses AES-256-CBC encryption algorithm
  - Supports parsing encMsgV3 fields from `info.xml`

## Algorithm Details

### checkMsgV3 (Signature Verification)

Format in `info.xml`: `${expectedHmac}${salt}_${filename}`, multiple items separated by `**`.

- **expectedHmac**: 64 characters
  - Calculation: `hmacSha256(input, lowercase(hexEncode(pbkdf2Key)) as []byte)`
  - pbkdf2Key: Generated using pbkdf2-sha256: `pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New)`
  - Note: pbkdf2Key needs hex encoding followed by lowercase conversion

- **salt**: 64 characters, used for pbkdf2

Example:
```
e56ac33a0eb3e97e501ded79eecc16496feb009a3ec46911186881f3dd73f3b7cec932efa6414914304a7e024f96686c38c7137bd734a407ba0a40d24696f813_com.tencent.mm514.tar**...
```

### encMsgV3 (AES Encryption Parameters)

96 characters in `info.xml`, format: `${salt}${iv}`:

- **salt**: First 64 characters (hex encoded)
- **iv**: Last 32 characters (hex encoded)

## Usage

### checkhash - Verify Backup Files

```sh
./checkhash \
  --password 12345678 \
  --checkMsgV3 50835ee73fb95dfe4712dd42ee926476887908d20e6d02c3800494f08dee77835e415a98c5553c85ff86446b61e753f5a62e7ed1dc45c072853f6e92e78bb283_com.tencent.mm0.tar \
  --input ./com.tencent.mm0.tar
```

Output example:
```
2024/10/10 00:33:49 checkMsgV3Item.ExpectedHmac: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 checkMsgV3Item.Salt: 5E415A98C5553C85FF86446B61E753F5A62E7ED1DC45C072853F6E92E78BB283
2024/10/10 00:33:49 checkMsgV3Item.FileName: com.tencent.mm0.tar
2024/10/10 00:33:49 File Hash: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 Success
```

### decrypt - Decrypt Backup Files

```sh
./decrypt \
  --password 12345678 \
  --encMsgV3 0ea2404230f7d824b354feea5d5cec6b24fe35303d4a9d9f687d0641aa5f19a3226264ab0ba258e1dca455d032d19de6 \
  --input ./com.tencent.mm0.tar \
  --output ./out
```

Output example:
```
2024/10/10 00:43:46 encMsgV3.Salt: 0EA2404230F7D824B354FEEA5D5CEC6B24FE35303D4A9D9F687D0641AA5F19A3
2024/10/10 00:43:46 encMsgV3.Iv: 226264AB0BA258E1DCA455D032D19DE6
2024/10/10 00:43:46 key: 85C973940B15BA7B7FBC203025D86B38B888EDD0CD3577F16ECE24BFC951962D
2024/10/10 00:43:48 Success
```

## Test Environment

- backupVersion: 29
- backupVersionName: 13.1.0.340
- Device: HMA-AL00
- Hisuite: 14.0.0.320_OVE

## References

- [Huawei-Hisuite-KobackupCipherTool](https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool)
- [KobackupCipherTool](https://github.com/oO0oO0oO0o0o00/KobackupCipherTool)
- [Huawei-Hisuite-KobackupCipherTool (ctwolfxzd)](https://github.com/ctwolfxzd/Huawei-Hisuite-KobackupCipherTool)
- [kobackupdec](https://github.com/RealityNet/kobackupdec)

## License

- No commercial use allowed
- For personal learning and educational purposes only
- No redistribution allowed
