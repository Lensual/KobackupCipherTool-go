# KobackupCipherTool-go

- checkMsgV3: 使用备份密码和 hmacSha256 算法对备份文件的签名，用于校验
  - 在 `info.xml` 中，格式为 `${expectedHmac}${salt}_${filename}` ，以 `**` 作为多个项目的分隔符，最后一个项目含义不明
    - expectedHmac: 64 个字符
      - 计算方法 `hmacSha256(input, lowercase(hexEncode(pbkdf2Key)) as []byte)`
      - pbkdf2Key: 对备份时的密码使用 pbkdf2-sha256 生成密钥 `pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New)`
      - 注意对 pbkdf2Key 进行 hex 编码，转换为小写
    - salt: 64 个字符，用于 pbkdf2
  - 示例: `e56ac33a0eb3e97e501ded79eecc16496feb009a3ec46911186881f3dd73f3b7cec932efa6414914304a7e024f96686c38c7137bd734a407ba0a40d24696f813_com.tencent.mm514.tar**360194b1f3ff3119f40dedd35b8155dce5896e159a79e62fbd0423f80bad54768d7b9d9f777d07797cfaf8ccce9142bd548c6a4b05e415f115215dfd82a61dfb_com.tencent.mm`

- encMsgV3: AES 加密参数
  - 在 `info.xml` 中，96 个字符，格式为 `${salt}${iv}`
  - salt: 前 64 个字符，hex 编码
  - iv: 后 32 个字符，hex 编码

## Usage

checkhash

```sh
./checkhash \
  --password 12345678 \
  --checkMsgV3 50835ee73fb95dfe4712dd42ee926476887908d20e6d02c3800494f08dee77835e415a98c5553c85ff86446b61e753f5a62e7ed1dc45c072853f6e92e78bb283_com.tencent.mm0.tar \
  --input ./com.tencent.mm0.tar

2024/10/10 00:33:49 checkMsgV3Item.ExpectedHmac: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 checkMsgV3Item.Salt: 5E415A98C5553C85FF86446B61E753F5A62E7ED1DC45C072853F6E92E78BB283
2024/10/10 00:33:49 checkMsgV3Item.FileName: com.tencent.mm0.tar
2024/10/10 00:33:49 File Hash: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 Success
```

```sh
./decrypt \
  --password 12345678 \
  --encMsgV3 0ea2404230f7d824b354feea5d5cec6b24fe35303d4a9d9f687d0641aa5f19a3226264ab0ba258e1dca455d032d19de6 \
  --input ./com.tencent.mm0.tar \
  --output ./out
                
2024/10/10 00:43:46 encMsgV3.Salt: 0EA2404230F7D824B354FEEA5D5CEC6B24FE35303D4A9D9F687D0641AA5F19A3
2024/10/10 00:43:46 encMsgV3.Iv: 226264AB0BA258E1DCA455D032D19DE6
2024/10/10 00:43:46 key: 85C973940B15BA7B7FBC203025D86B38B888EDD0CD3577F16ECE24BFC951962D
2024/10/10 00:43:48 Success
```

## Others

Tested on `backupVersion 29; backupVersionName 13.1.0.340; HMA-AL00; Hisuite 14.0.0.320_OVE`

## Refrences

[https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool](https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool)

[https://github.com/oO0oO0oO0o0o00/KobackupCipherTool](https://github.com/oO0oO0oO0o0o00/KobackupCipherTool)

[https://github.com/ctwolfxzd/Huawei-Hisuite-KobackupCipherTool](https://github.com/ctwolfxzd/Huawei-Hisuite-KobackupCipherTool)

[https://github.com/RealityNet/kobackupdec](https://github.com/RealityNet/kobackupdec)

## LICENSE

- 禁止用于任何商业目的
- 仅限个人学习、教学用途
- 禁止二次分发
