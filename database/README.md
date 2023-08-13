## 百科格式

```yml
# 該東西的網址跟代號，在其他地方提及時會用到。
# 例如：`asus` 是華碩公司，如果是華碩的子公司則需以 `asus_` 為前輟。
code: asus

# 上級代號（可選）。
# 例如：`asus_rog` 玩家國度的上級代號應該填入 `asus` 華碩公司。
grand_code: ""

# 英文全名（可選）。
fullname_en: ASUSTeK Computer (Asus)

# 中文全名（可選）。
fullname_zh: 華碩電腦

# 其他稱呼（可選）。
subnames: [華碩, ASUS]

# 官方網站或是社群網站網址（可選）。
website: https://asus.com

# 公司簡介（可選）。
description: Asus（TWSE：2357）總部位於臺灣。

# 物件的分類，請參閱 `README_CATEGORIES.md` 檔案，可複選。
categories: [IT, MP, PC]

# 該物件的性質：
# brand（品牌、組織）；company（立案公司）；software（軟體、網站、服務）；game（遊戲）
type: company

# 哪個國家擁有：
# 台灣：taiwanese_gov（政府）；taiwanese（公司）；taiwanese_fake（空殼公司）
# 香港：chinese_gov（政府）；hongkongese（公司）；hongkongese_fake（空殼公司）
# 中國：chinese_gov（政府）；chinese（公司）
# 外國：foreign_gov（政府）；foreign（公司）；foreign_fake（空殼公司）
#
# 空殼公司主要目的是為了合法化洗刷中資來歷，例如：
# -《中國 Bilibili》在台灣註冊《小萌科技》並偽裝是遊戲新創，即是 taiwanese_fake。
# -《中國 小米》成立的《華米穿戴品牌》表面上是國際公司，但技術核心來自中國，即是 foreign_fake。
owner: taiwanese

# 歷史記事。
informations:
    # 日期（YYYY-MM-DD 或 YYYY），留空表示不特別記載。
    - date: 2020-09-24

      # 說明
      description: |
          相較於其他國家，華碩在中國有更多品牌且提供更多客製化主機板與電腦硬體。

      # 收入來源變更（可選）：
      # 中國政府：chinese_gov（由中國政府直接營運或管控）
      # 中國股東：chinese_shareholder（控制公司、組織的股份是中國人，且已經過中國相關法律核准）
      # 中國市場：chinese（大多數的收入來源在中國、中國被視為不可放棄的市場）
      # 不再追蹤：untracked（收入來源已經不在中國）
      income: chinese_shareholder

      # 營運方式變更（可選）：
      # 跟隨中國：chinese_gov（跟隨中國政府的主要政策）
      # 討好中國：chinese（透過符合中國人民喜好來從中獲得利益）
      # 試圖中立：neutral（避免任何政治問題並試圖達到政治中立）
      # 不再追蹤：untracked（已經不再受中國控制）
      direction: untracked

      # 支持侵略變更（可選）：
      # 支持侵略：support（表明支持中國侵略、統一；將自己視為中國持有而認同中國的政策）
      # 間接協助：indirect（未表明支持，但透過間接行為協助中國助長）
      # 未表態　：neutral（明顯地避而不談）
      # 不再追蹤：untracked（不傾向中國）
      invasion: untracked
```

## 詞彙格式

```yml
# 詞彙本身。
word: 計算機

# 注音，每個單字以空白隔開。
bopomofo: ㄐㄧˋ ㄙㄨㄢˋ ㄐㄧ

# 說明（可選）。
description: |
    個人或是辦公用的電腦，被用作娛樂或是各種行業的數位裝置，手機也算是可攜式電腦。

# 區分說明（可選）。
distinguish: |
    隨手可得的小型數學計算機

# 範例。
examples:
    # 正確的詞彙，同個範例可以有多個正確的詞彙。
    - words: [電腦]

      # 說明（可選）。
      # 如果這個說明跟詞彙說明一樣，通常建議挪到詞彙說明即可。
      description: |
          個人或是辦公用的電腦，被用作娛樂或是各種行業的數位裝置，手機也算是可攜式電腦。

      # 正確的範例（可選）。
      # 使用 {{文字}} 來標記詞彙，而 [[文字]] 會產生連結到其他詞彙。
      correct: |
          我剛買了最新的{{電腦}}。

      # 不正確的範例（可選）。
      # 使用 {{文字}} 來標記詞彙，而 [[文字]] 會產生連結到其他詞彙。
      incorrect: |
          我剛買了最新的{{計算機}}。
```

## 分類

```
GP：集團
FI：銀行/金融
PP：房地產/能源
PL：政治
ME：傳播媒體
BK：教育/書商/出版社
SP：購物
CY：日常用品
EA：家具/電器
MD：醫藥/美容
FD：飲食
ST：運動
EN：娛樂
GA：遊戲
AP：軟體/網頁
IT：資訊與科技
MP：通訊商/手機
PC：電腦/硬體
TS：交通
TL：旅行
```
