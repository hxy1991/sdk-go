package utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestReplaceAllPunctuation(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				s: "*",
			},
			want: "*",
		},
		{
			name: "",
			args: args{
				s: "你好，世界！",
			},
			want: "你好*世界*",
		},
		{
			name: "",
			args: args{
				s: "Hello,world!",
			},
			want: "Hello*world*",
		},
		{
			name: "",
			args: args{
				s: "Hello,world!",
			},
			want: "Hello*world*",
		},
		{
			name: "",
			args: args{
				s: "*ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
			},
			want: "*ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		},
		{
			name: "",
			args: args{
				// 中文标点符号
				s: "0。0？0！0，0、0；0：0“0”0‘0’「0」0『0』0（0）0[0]0〔0〕0【0】0—0—0…0…0·0—0-0~0·0《0》0〈0〉0_0_0_0/0",
			},
			want: "0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0",
		},
		{
			name: "",
			args: args{
				// 英文标点符号
				s: "0,0:0;0‐0’0'0′0″0‴0.0&0@0^0/0\\0…0*0⁂0***0-0-0‒0–0—0=0⸗0?0!0‽0¡0¿0!?0⸮0#0№0º0ª0%0‰0‱0°0⌀0+0−0×0÷0~0±0∓0–0_0⁀0|0¦0‖0•0·0©0©0℗0®0‘’0“”0''0\"\"0‹›0«»0()0[]0{}0⟨⟩0”0〃0†0‡0❧0/0☞0◊0¶0⸿0፠0๛0()0※0§0",
			},
			// ºª
			want: "0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0º0ª0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0*0",
		},
		{
			name: "",
			//https://cn.piliapp.com/symbol/
			//http://blog.neten.de/posts/333/
			args: args{
				s: "㏠ ㏡ ㏢ ㏣ ㏤ ㏥ ㏦ ㏧ ㏨ ㏩ ㏪ ㏫ ㏬ ㏭ ㏮ ㏯ ㏰ ㏱ ㏲ ㏳ ㏴ ㏵ ㏶ ㏷ ㏸ ㏹ ㏺ ㏻ ㏼ ㏽ ㏾㍘ ㍙ ㍚ ㍛ ㍜ ㍝ ㍞ ㍟ ㍠ ㍡ ㍢ ㍣ ㍤ ㍥ ㍦ ㍧ ㍨ ㍩ ㍪ ㍫ ㍬ ㍭ ㍮ ㍯ ㍰㋀ ㋁ ㋂ ㋃ ㋄ ㋅ ㋆ ㋇ ㋈ ㋉ ㋊ ㋋㉿ ㊀ ㊁ ㊂ ㊃ ㊄ ㊅ ㊆ ㊇ ㊈ ㊉ ㊊ ㊋ ㊌ ㊍ ㊎ ㊏ ㊐ ㊑ ㊒ ㊓ ㊔ ㊕ ㊖ ㊗ ㊘ ㊙ ㊚ ㊛ ㊜ ㊝ ㊞ ㊟ ㊠ ㊡ ㊢ ㊣ ㊤ ㊥ ㊦ ㊧ ㊨ ㊩ ㊪ ㊫ ㊬ ㊭ ㊮ ㊯ ㊰㈠ ㈡ ㈢ ㈣ ㈤ ㈥ ㈦ ㈧ ㈨ ㈩ ㈪ ㈫ ㈬ ㈭ ㈮ ㈯ ㈰ ㈱ ㈲ ㈳ ㈴ ㈵ ㈶ ㈷ ㈸ ㈹ ㈺ ㈻ ㈼ ㈽ ㈾ ㈿ ㉀ ㉁ ㉂ ㉃㆒ ㆓ ㆔ ㆕ ㆖ ㆗ ㆘ ㆙ ㆚ ㆛ ㆜ ㆝ ㆞ ㆟ ✀ ✁ ✂ ✃ ✄ ✅ ✆ ✇ ✈ ✉ ✊ ✋ ✌ ✍ ✎ ✏ ✐ ✑ ✒ ✓ ✔ ✕ ✖ ✗ ✘ ✙ ✚ ✛ ✜ ✝ ✞ ✟ ✠ ✡ ✢ ✣ ✤ ✥ ✦ ✧ ✨ ✩ ✪ ✫ ✬ ✭ ✮ ✯ ✰ ✱ ✲ ✳ ✴ ✵ ✶ ✷ ✸ ✹ ✺ ✻ ✼ ✽ ✾ ✿ ❀ ❁ ❂ ❃ ❄ ❅ ❆ ❇ ❈ ❉ ❊ ❋ ❌ ❍ ❎ ❏ ❐ ❑ ❒ ❓ ❔ ❕ ❖ ❗ ❘ ❙ ❚ ❛ ❜ ❝ ❞ ❟ ❠ ❡ ❢ ❣ ❤ ❥ ❦ ❧ ❨ ❩ ❪ ❫ ❬ ❭ ❮ ❯ ❰ ❱ ❲ ❳ ❴ ❵ ❶ ❷ ❸ ❹ ❺ ❻ ❼ ❽ ❾ ❿ ➀ ➁ ➂ ➃ ➄ ➅ ➆ ➇ ➈ ➉ ➊ ➋ ➌ ➍ ➎ ➏ ➐ ➑ ➒ ➓ ➔ ➕ ➖ ➗ ➘ ➙ ➚ ➛ ➜ ➝ ➞ ➟ ➠ ➡ ➢ ➣ ➤ ➥ ➦ ➧ ➨ ➩ ➪ ➫ ➬ ➭ ➮ ➯ ➰ ➱ ➲ ➳ ➴ ➵ ➶ ➷ ➸ ➹ ➺ ➻ ➼ ➽ ➾ ➿ ① ② ③ ④ ⑤ ⑥ ⑦ ⑧ ⑨ ⑩ ⑪ ⑫ ⑬ ⑭ ⑮ ⑯ ⑰ ⑱ ⑲ ⑳ ⑴ ⑵ ⑶ ⑷ ⑸ ⑹ ⑺ ⑻ ⑼ ⑽ ⑾ ⑿ ⒀ ⒁ ⒂ ⒃ ⒄ ⒅ ⒆ ⒇ ⒈ ⒉ ⒊ ⒋ ⒌ ⒍ ⒎ ⒏ ⒐ ⒑ ⒒ ⒓ ⒔ ⒕ ⒖ ⒗ ⒘ ⒙ ⒚ ⒛ ⒜ ⒝ ⒞ ⒟ ⒠ ⒡ ⒢ ⒣ ⒤ ⒥ ⒦ ⒧ ⒨ ⒩ ⒪ ⒫ ⒬ ⒭ ⒮ ⒯ ⒰ ⒱ ⒲ ⒳ ⒴ ⒵ Ⓐ Ⓑ Ⓒ Ⓓ Ⓔ Ⓕ Ⓖ Ⓗ Ⓘ Ⓙ Ⓚ Ⓛ Ⓜ Ⓝ Ⓞ Ⓟ Ⓠ Ⓡ Ⓢ Ⓣ Ⓤ Ⓥ Ⓦ Ⓧ Ⓨ Ⓩ ⓐ ⓑ ⓒ ⓓ ⓔ ⓕ ⓖ ⓗ ⓘ ⓙ ⓚ ⓛ ⓜ ⓝ ⓞ ⓟ ⓠ ⓡ ⓢ ⓣ ⓤ ⓥ ⓦ ⓧ ⓨ ⓩ ⓪ ⓫ ⓬ ⓭ ⓮ ⓯ ⓰ ⓱ ⓲ ⓳ ⓴ ⓵ ⓶ ⓷ ⓸ ⓹ ⓺ ⓻ ⓼ ⓽ ⓾ ⓿ ─ ━ │ ┃ ┄ ┅ ┆ ┇ ┈ ┉ ┊ ┋ ┌ ┍ ┎ ┏ ┐ ┑ ┒ ┓ └ ┕ ┖ ┗ ┘ ┙ ┚ ┛ ├ ┝ ┞ ┟ ┠ ┡ ┢ ┣ ┤ ┥ ┦ ┧ ┨ ┩ ┪ ┫ ┬ ┭ ┮ ┯ ┰ ┱ ┲ ┳ ┴ ┵ ┶ ┷ ┸ ┹ ┺ ┻ ┼ ┽ ┾ ┿ ╀ ╁ ╂ ╃ ╄ ╅ ╆ ╇ ╈ ╉ ╊ ╋ ╌ ╍ ╎ ╏ ═ ║ ╒ ╓ ╔ ╕ ╖ ╗ ╘ ╙ ╚ ╛ ╜ ╝ ╞ ╟ ╠ ╡ ╢ ╣ ╤ ╥ ╦ ╧ ╨ ╩ ╪ ╫ ╬ ╭ ╮ ╯ ╰ ╱ ╲ ╳ ╴ ╵ ╶ ╷ ╸ ╹ ╺ ╻ ╼ ╽ ╾ ╿ ▀ ▁ ▂ ▃ ▄ ▅ ▆ ▇ █ ▉ ▊ ▋ ▌ ▍ ▎ ▏ ▐ ░ ▒ ▓ ▔ ▕ ▖ ▗ ▘ ▙ ▚ ▛ ▜ ▝ ▞ ▟ ■ □ ▢ ▣ ▤ ▥ ▦ ▧ ▨ ▩ ▪ ▫ ▬ ▭ ▮ ▯ ▰ ▱ ▲ △ ▴ ▵ ▶ ▷ ▸ ▹ ► ▻ ▼ ▽ ▾ ▿ ◀ ◁ ◂ ◃ ◄ ◅ ◆ ◇ ◈ ◉ ◊ ○ ◌ ◍ ◎ ● ◐ ◑ ◒ ◓ ◔ ◕ ◖ ◗ ◘ ◙ ◚ ◛ ◜ ◝ ◞ ◟ ◠ ◡ ◢ ◣ ◤ ◥ ◦ ◧ ◨ ◩ ◪ ◫ ◬ ◭ ◮ ◯ ◰ ◱ ◲ ◳ ◴ ◵ ◶ ◷ ◸ ◹ ◺ ◻ ◼ ◽ ◾ ◿ ☀ ☁ ☂ ☃ ☄ ★ ☆ ☇ ☈ ☉ ☊ ☋ ☌ ☍ ☎ ☏ ☐ ☑ ☒ ☓ ☔ ☕ ☖ ☗ ☘ ☙ ☚ ☛ ☜ ☝ ☞ ☟ ☠ ☡ ☢ ☣ ☤ ☥ ☦ ☧ ☨ ☩ ☪ ☫ ☬ ☭ ☮ ☯ ☰ ☱ ☲ ☳ ☴ ☵ ☶ ☷ ☸ ☹ ☺ ☻ ☼ ☽ ☾ ☿ ♀ ♁ ♂ ♃ ♄ ♅ ♆ ♇ ♈ ♉ ♊ ♋ ♌ ♍ ♎ ♏ ♐ ♑ ♒ ♓ ♔ ♕ ♖ ♗ ♘ ♙ ♚ ♛ ♜ ♝ ♞ ♟ ♠ ♡ ♢ ♣ ♤ ♥ ♦ ♧ ♨ ♩ ⌠ ⌡ ⌢ ⌣⌆ ⌇ ⌈ ⌉ ⌊ ⌋ ∀ ∁ ∂ ∃ ∄ ∅ ∆ ∇ ∈ ∉ ∊ ∋ ∌ ∍ ∎ ∏ ∐ ∑ − ∓ ∔ ∕ ∖ ∗ ∘ ∙ √ ∛ ∜ ∝ ∞ ∟ ∠ ∡ ∢ ∣ ∤ ∥ ∦ ∧ ∨ ∩ ∪ ∫ ∬ ∭ ∮ ∯ ∰ ∱ ∲ ∳ ∴ ∵ ∶ ∷ ∸ ∹ ∺ ∻ ∼ ∽ ∾ ∿ ≀ ≁ ≂ ≃ ≄ ≅ ≆ ≇ ≈ ≉ ≊ ≋ ≌ ≍ ≎ ≏ ≐ ≑ ≒ ≓ ≔ ≕ ≖ ≗ ≘ ≙ ≚ ≛ ≜ ≝ ≞ ≟ ≠ ≡ ≢ ≣ ≤ ≥ ≦ ≧ ≨ ≩ ≪ ≫ ≬ ≭ ≮ ≯ ≰ ≱ ≲ ≳ ≴ ≵ ≶ ≷ ≸ ≹ ≺ ≻ ≼ ≽ ≾ ≿ ⊀ ⊁ ⊂ ⊃ ⊄ ⊅ ⊆ ⊇ ⊈ ⊉ ⊊ ⊋ ⊌ ⊍ ⊎ ⊏ ⊐ ⊑ ⊒ ⊓ ⊔ ⊕ ⊖ ⊗ ⊘ ⊙ ⊚ ⊛ ⊜ ⊝ ⊞ ⊟ ⊠ ⊡ ⊢ ⊣ ⊤ ⊥ ⊦ ⊧ ⊨ ⊩ ⊪ ⊫ ⊬ ⊭ ⊮ ⊯ ⊰ ⊱ ⊲ ⊳ ⊴ ⊵ ⊶ ⊷ ⊸ ⊹ ⊺ ⊻ ⊼ ⊽ ⊾ ⊿ ⋀ ⋁ ⋂ ⋃ ⋄ ⋅ ⋆ ⋇ ⋈ ⋉ ⋊ ⋋ ⋌ ⋍ ⋎ ⋏ ⋐ ⋑ ⋒ ⋓ ⋔ ⋕ ⋖ ⋗ ⋘ ⋙ ⋚ ⋛ ⋜ ⋝ ⋞ ⋟ ⋠ ⋡ ⋢ ⋣ ⋤ ⋥ ⋦ ⋧ ⋨ ⋩ ⋪ ⋫ ⋬ ⋭ ⋮ ⋯ ⋰ ⋱ ← ↑ → ↓ ↔ ↕ ↖ ↗ ↘ ↙ ↚ ↛ ↜ ↝ ↞ ↟ ↠ ↡ ↢ ↣ ↤ ↥ ↦ ↧ ↨ ↩ ↪ ↫ ↬ ↭ ↮ ↯ ↰ ↱ ↲ ↳ ↴ ↵ ↶ ↷ ↸ ↹ ↺ ↻ ↼ ↽ ↾ ↿ ⇀ ⇁ ⇂ ⇃ ⇄ ⇅ ⇆ ⇇ ⇈ ⇉ ⇊ ⇋ ⇌ ⇍ ⇎ ⇏ ⇐ ⇑ ⇒ ⇓ ⇔ ⇕ ⇖ ⇗ ⇘ ⇙ ⇚ ⇛ ⇜ ⇝ ⇞ ⇟ ⇠ ⇡ ⇢ ⇣ ⇤ ⇥ ⇦ ⇧ ⇨ ⇩ ⇪⅓ ⅔ ⅕ ⅖ ⅗ ⅘ ⅙ ⅚ ⅛ ⅜ ⅝ ⅞ ⅟ Ⅰ Ⅱ Ⅲ Ⅳ Ⅴ Ⅵ Ⅶ Ⅷ Ⅸ Ⅹ Ⅺ Ⅻ Ⅼ Ⅽ Ⅾ Ⅿ ⅰ ⅱ ⅲ ⅳ ⅴ ⅵ ⅶ ⅷ ⅸ ⅹ ⅺ ⅻ ⅼ ⅽ ⅾ ⅿ ↀ ↁ ↂ℀ ℁ ℂ ℃ ℄ ℅ ℆ ℇ ℈ ℉ ℊ ℋ ℌ ℍ ℎ ℏ ℐ ℑ ℒ ℓ ℔ ℕ № ℗ ℘ ℙ ℚ ℛ ℜ ℝ ℞ ℟ ℠ ℡ ™ ℣ ℤ ℥ Ω ℧ ℨ ℩ K Å ℬ ℭ ℮ ℯ ℰ ℱ Ⅎ ℳ ℴ ℵ ℶ ℷ ℸ ₠ ₡ ₢ ₣ ₤ ₥ ₦ ₧ ₨ ₩ ₪ ₫ € ₭ ₮ ₯ ⁺ ⁻ ⁼ ⁽ ⁾ ⁿ ₀ ₁ ₂ ₃ ₄ ₅ ₆ ₇ ₈ ₉ ₊ ₋ ₌ ₍ ₎‐ ‑ ‒ – — ― ‖ ‗ ‘ ’ ‚ ‛ “ ” „ ‟ † ‡ • ‣ ․ ‥ … ‧  ",
			},
			// ㊀㊁㊂㊃㊄㊅㊆㊇㊈㊉㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩㆒㆓㆔㆕❶❷❸❹❺❻❼❽❾❿➀➁➂➃➄➅➆➇➈➉➊➋➌➍➎➏➐➑➒➓①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛⓪⓫⓬⓭⓮⓯⓰⓱⓲⓳⓴⓵⓶⓷⓸⓹⓺⓻⓼⓽⓾⓿⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞⅟ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫⅬⅭⅮⅯⅰⅱⅲⅳⅴⅵⅶⅷⅸⅹⅺⅻⅼⅽⅾⅿↀↁↂℂℇℊℋℌℍℎℏℐℑℒℓℕℙℚℛℜℝℤΩℨKÅℬℭℯℰℱℲℳℴℵℶℷℸⁿ₀₁₂₃₄₅₆₇₈₉ ,
			want: "* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * ㊀ ㊁ ㊂ ㊃ ㊄ ㊅ ㊆ ㊇ ㊈ ㊉ * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *㈠ ㈡ ㈢ ㈣ ㈤ ㈥ ㈦ ㈧ ㈨ ㈩ * * * * * * * * * * * * * * * * * * * * * * * * * *㆒ ㆓ ㆔ ㆕ * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * ❶ ❷ ❸ ❹ ❺ ❻ ❼ ❽ ❾ ❿ ➀ ➁ ➂ ➃ ➄ ➅ ➆ ➇ ➈ ➉ ➊ ➋ ➌ ➍ ➎ ➏ ➐ ➑ ➒ ➓ * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * ① ② ③ ④ ⑤ ⑥ ⑦ ⑧ ⑨ ⑩ ⑪ ⑫ ⑬ ⑭ ⑮ ⑯ ⑰ ⑱ ⑲ ⑳ ⑴ ⑵ ⑶ ⑷ ⑸ ⑹ ⑺ ⑻ ⑼ ⑽ ⑾ ⑿ ⒀ ⒁ ⒂ ⒃ ⒄ ⒅ ⒆ ⒇ ⒈ ⒉ ⒊ ⒋ ⒌ ⒍ ⒎ ⒏ ⒐ ⒑ ⒒ ⒓ ⒔ ⒕ ⒖ ⒗ ⒘ ⒙ ⒚ ⒛ * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * ⓪ ⓫ ⓬ ⓭ ⓮ ⓯ ⓰ ⓱ ⓲ ⓳ ⓴ ⓵ ⓶ ⓷ ⓸ ⓹ ⓺ ⓻ ⓼ ⓽ ⓾ ⓿ * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *⅓ ⅔ ⅕ ⅖ ⅗ ⅘ ⅙ ⅚ ⅛ ⅜ ⅝ ⅞ ⅟ Ⅰ Ⅱ Ⅲ Ⅳ Ⅴ Ⅵ Ⅶ Ⅷ Ⅸ Ⅹ Ⅺ Ⅻ Ⅼ Ⅽ Ⅾ Ⅿ ⅰ ⅱ ⅲ ⅳ ⅴ ⅵ ⅶ ⅷ ⅸ ⅹ ⅺ ⅻ ⅼ ⅽ ⅾ ⅿ ↀ ↁ ↂ* * ℂ * * * * ℇ * * ℊ ℋ ℌ ℍ ℎ ℏ ℐ ℑ ℒ ℓ * ℕ * * * ℙ ℚ ℛ ℜ ℝ * * * * * * ℤ * Ω * ℨ * K Å ℬ ℭ * ℯ ℰ ℱ Ⅎ ℳ ℴ ℵ ℶ ℷ ℸ * * * * * * * * * * * * * * * * * * * * * ⁿ ₀ ₁ ₂ ₃ ₄ ₅ ₆ ₇ ₈ ₉ * * * * * * * * * * * * * * * * * * * * * * * * * * * *  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containPunctuation := ContainPunctuationOrSymbol(tt.args.s)
			assert.True(t, containPunctuation)

			if got := ReplaceAllPunctuationOrSymbol(tt.args.s, "*"); got != tt.want {
				t.Errorf("ReplaceAllPunctuation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceAllEmoji(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				s: "😀😃😄😁😆😅🤣😂🙂🙃",
			},
			want: "**********",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containPunctuation := ContainEmoji(tt.args.s)
			assert.True(t, containPunctuation)

			if got := ReplaceAllEmoji(tt.args.s, "*"); got != tt.want {
				t.Errorf("ReplaceAllEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceWhitespace(t *testing.T) {
	type args struct {
		src  string
		repl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				src:  "0 0\t0\r0\n0",
				repl: "*",
			},
			want: "0*0*0*0*0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containWhitespace := ContainWhitespace(tt.args.src)
			assert.True(t, containWhitespace)

			if got := ReplaceWhitespace(tt.args.src, tt.args.repl); got != tt.want {
				t.Errorf("ReplaceWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContain(t *testing.T) {
	type args struct {
		reg *regexp.Regexp
		src string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				reg: regexp.MustCompile("^\\p{L}*$"),
				src: "中abc",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contain(tt.args.reg, tt.args.src); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnlyLanguageCharacters(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				// 简体常用字
				// https://baike.baidu.com/item/%E7%8E%B0%E4%BB%A3%E6%B1%89%E8%AF%AD%E5%B8%B8%E7%94%A8%E5%AD%97%E8%A1%A8/8922402?fromModule=lemma_inlink
				"一乙二十丁厂七卜人入八九几儿了力乃刀又三于干亏士工土才寸下大丈与万上小口巾山千乞川亿个勺久凡及夕丸么广亡门义之尸弓己已子",
			},
			want: true,
		},
		{
			// 台湾繁体常用字
			// https://zh.wikipedia.org/zh-sg/%E5%B8%B8%E7%94%A8%E5%9C%8B%E5%AD%97%E6%A8%99%E6%BA%96%E5%AD%97%E9%AB%94%E8%A1%A8
			// https://web.archive.org/web/20210514103801/https://language.moe.gov.tw/001/Upload/Files/site_content/download/mandr/%E6%95%99%E8%82%B2%E9%83%A84808%E5%80%8B%E5%B8%B8%E7%94%A8%E5%AD%97.pdf
			name: "",
			args: args{
				"一丁七三下丈上丑丐不丙世丕且丘丞丟並丫中串丸凡丹主乃久么之尹乍乏乎乒乓乖乘乙九也乞乩乳乾亂了予事二于云井互五亙些亞亟亡交亦",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				"0一",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				"0😀一",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				"0㏠一",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				"0 一",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				"0!一",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				"0！一",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				"Hello0Imonenicetomeetyou",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 日
				"こんにちは0はじめまして",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 韩
				"안녕하세요0만나서반가워요",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 西班牙
				"bonjour0jesuisravidevousrencontrer",
			},
			want: true,
		},
		//{
		//	name: "",
		//	args: args{
		//		//   TODO 葡萄牙会用到-，如 conhecê-lo 认识你
		//		"Olá0eusouumprazerconhecê-lo",
		//	},
		//	want: false,
		//},
		{
			name: "",
			args: args{
				// 俄
				"привет0ярадпознакомитьсясвами",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 越南
				"xinchào0tôirấtvuiđượcgặpbạn",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 泰国 \p{M}
				"ฉันดีใจที่ได้พบคุณ",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 土耳其
				"merhabatanıştığımamemnunoldum",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 阿拉伯
				"مرحبااناسعيدبلقائك",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OnlyLanguageOrNumberCharacters(tt.args.src); got != tt.want {
				t.Errorf("OnlyLanguageOrNumberCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainChinese(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				"Hello0Imonenicetomeetyou",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				// 日
				"こんにちは0はじめまして",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				// 韩
				"안녕하세요0만나서반가워요",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				// 西班牙
				"bonjour0jesuisravidevousrencontrer",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				// 中文
				"a你好0我很高兴认识你b",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				// 中文
				"1你好，我很高兴认识你",
			},
			want: true,
		},

		{
			name: "",
			args: args{
				// 中文
				"你好，我很高兴认识你",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainChinese(tt.args.src), "ContainChinese(%v)", tt.args.src)
		})
	}
}
