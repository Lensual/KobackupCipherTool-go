# KobackupCipherTool-go

Kobackup 备份文件解密工具，使用 Go 语言实现。用于解析和解密 Kobackup 的加密文件。

## 功能特性

- **checkhash**: 验证备份文件完整性
  - 使用备份密码和 HMAC-SHA256 算法验证签名
  - 支持解析 `info.xml` 中的 checkMsgV3 字段

- **decrypt**: 解密备份文件
  - 使用 AES-256-GCM 加密算法
  - 支持解析 `info.xml` 中的 encMsgV3 字段

- **decrypt-dir**: 批量解密整个备份目录
  - 自动从 `info.xml` 解析加密参数
  - 自动从 `backupinfo.ini` 获取应用包名
  - 递归解密目录下所有 `.tar` 文件

## 算法说明

### checkMsgV3（签名验证）

在 `info.xml` 中格式为 `${expectedHmac}${salt}_${filename}`，多个项目以 `**` 分隔。

- **expectedHmac**: 64 个字符
  - 计算方法：`hmacSha256(input, lowercase(hexEncode(pbkdf2Key)) as []byte)`
  - pbkdf2Key：使用 pbkdf2-sha256 生成 `pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New)`
  - 注意：pbkdf2Key 需要 hex 编码后转换为小写

- **salt**: 64 个字符，用于 pbkdf2

示例：
```
e56ac33a0eb3e97e501ded79eecc16496feb009a3ec46911186881f3dd73f3b7cec932efa6414914304a7e024f96686c38c7137bd734a407ba0a40d24696f813_com.tencent.mm514.tar**...
```

### encMsgV3（AES 加密参数）

在 `info.xml` 中为 96 个字符，格式为 `${salt}${iv}`：

- **salt**: 前 64 个字符（hex 编码）
- **iv**: 后 32 个字符（hex 编码）

## 使用方法

### checkhash - 验证备份文件

```sh
./checkhash \
  --password 12345678 \
  --checkMsgV3 50835ee73fb95dfe4712dd42ee926476887908d20e6d02c3800494f08dee77835e415a98c5553c85ff86446b61e753f5a62e7ed1dc45c072853f6e92e78bb283_com.tencent.mm0.tar \
  --input ./com.tencent.mm0.tar
```

输出示例：
```
2024/10/10 00:33:49 checkMsgV3Item.ExpectedHmac: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 checkMsgV3Item.Salt: 5E415A98C5553C85FF86446B61E753F5A62E7ED1DC45C072853F6E92E78BB283
2024/10/10 00:33:49 checkMsgV3Item.FileName: com.tencent.mm0.tar
2024/10/10 00:33:49 File Hash: 50835EE73FB95DFE4712DD42EE926476887908D20E6D02C3800494F08DEE7783
2024/10/10 00:33:49 Success
```

### decrypt - 解密备份文件

```sh
./decrypt \
  --password 12345678 \
  --encMsgV3 0ea2404230f7d824b354feea5d5cec6b24fe35303d4a9d9f687d0641aa5f19a3226264ab0ba258e1dca455d032d19de6 \
  --input ./com.tencent.mm0.tar \
  --output ./out

输出示例：
```
2024/10/10 00:43:46 encMsgV3.Salt: 0EA2404230F7D824B354FEEA5D5CEC6B24FE35303D4A9D9F687D0641AA5F19A3
2024/10/10 00:43:46 encMsgV3.Iv: 226264AB0BA258E1DCA455D032D19DE6
2024/10/10 00:43:46 key: 85C973940B15BA7B7FBC203025D86B38B888EDD0CD3577F16ECE24BFC951962D
2024/10/10 00:43:48 Success
```

### decrypt-dir - 批量解密备份目录

自动从目录中的 `info.xml` 和 `backupinfo.ini` 解析加密参数和包名，解密整个备份目录。

```sh
./decrypt-dir \
  --password 12345678 \
  --input ./backup_files

输出示例：
```
2024/10/10 00:43:46 encMsgV3 from info.xml: 0ea2404230f7d824b354feea5d5cec6b24fe35303d4a9d9f687d0641aa5f19a3226264ab0ba258e1dca455d032d19de6
2024/10/10 00:43:46 encMsgV3.Salt: 0EA2404230F7D824B354FEEA5D5CEC6B24FE35303D4A9D9F687D0641AA5F19A3
2024/10/10 00:43:46 encMsgV3.Iv: 226264AB0BA258E1DCA455D032D19DE6
2024/10/10 00:43:46 key: 85C973940B15BA7B7FBC203025D86B38B888EDD0CD3577F16ECE24BFC951962D
2024/10/10 00:43:46 Found package name: com.tencent.mm, targeting directory: ./backup_files/com.tencent.mm_appDataTar
2024/10/10 00:43:48 Decrypting: ./backup_files/com.tencent.mm_appDataTar/xxx.tar -> ./backup_files_decrypted/com.tencent.mm_appDataTar/xxx.tar
2024/10/10 00:43:48 Success: ./backup_files_decrypted/com.tencent.mm_appDataTar/xxx.tar
2024/10/10 00:43:49 Folder decryption completed
```

## 测试环境

成功

- backupVersion: 29
- backupVersionName: 13.1.0.340
- Device: HMA-AL00
- hisuiteversion: 14.0.0.320

失败

- backupVersion: 29
- backupVersionName: 14.5.0.375
- Device: HMA-AL00
- hisuiteversion: 14.0.0.340

## 参考项目

- [Huawei-Hisuite-KobackupCipherTool](https://github.com/irsl/Huawei-Hisuite-KobackupCipherTool)
- [KobackupCipherTool](https://github.com/oO0oO0oO0o0o00/KobackupCipherTool)
- [Huawei-Hisuite-KobackupCipherTool (ctwolfxzd)](https://github.com/ctwolfxzd/Huawei-Hisuite-KobackupCipherTool)
- [kobackupdec](https://github.com/RealityNet/kobackupdec)

## 许可证

- 禁止用于任何商业目的
- 仅限个人学习、教学用途
- 禁止二次分发
