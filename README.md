# KobackupCipherTool-go

WIP.

- checkMsgV3: 在 `info.xml` 中，格式为 `${expectedHmac}${salt}_${filename}` ，以 `**` 作为多个项目的分隔符，最后一个结尾没有分隔符
  - expectedHmac: 64 个字符，计算方法 `hmacSha256(input, hexEncode(hmacKey))`
    - hmacKey: 备份时的密码，使用 pbkdf2-sha256 生成密钥，iter 5000，keyLen 32
  - salt: 64 个字符
  - 示例：`e56ac33a0eb3e97e501ded79eecc16496feb009a3ec46911186881f3dd73f3b7cec932efa6414914304a7e024f96686c38c7137bd734a407ba0a40d24696f813_com.tencent.mm514.tar**360194b1f3ff3119f40dedd35b8155dce5896e159a79e62fbd0423f80bad54768d7b9d9f777d07797cfaf8ccce9142bd548c6a4b05e415f115215dfd82a61dfb_com.tencent.mm`

## Refrences

[https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool](https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool)

[https://github.com/RealityNet/kobackupdec](https://github.com/RealityNet/kobackupdec)
