#!/bin/bash

# RFQ 簡報轉換腳本
# 一鍵轉換 Markdown 為 PDF/HTML/PPTX

set -e

echo "🚀 RFQ 簡報轉換工具"
echo "===================="
echo ""

# 檢查 Marp CLI 是否已安裝
if ! command -v marp &> /dev/null; then
    echo "❌ Marp CLI 未安裝"
    echo ""
    echo "請執行以下指令安裝："
    echo "  npm install -g @marp-team/marp-cli"
    echo ""
    echo "或使用 Homebrew (macOS):"
    echo "  brew install marp-cli"
    echo ""
    exit 1
fi

echo "✅ Marp CLI 已安裝"
echo ""

# 切換到 rfq 目錄
cd "$(dirname "$0")"

# 檢查輸入檔案
INPUT_FILE="RFQ簡報.md"
if [ ! -f "$INPUT_FILE" ]; then
    echo "❌ 找不到檔案: $INPUT_FILE"
    exit 1
fi

echo "📄 輸入檔案: $INPUT_FILE"
echo ""

# 選單
echo "請選擇轉換格式："
echo "  1) PDF (建議用於正式簡報)"
echo "  2) HTML (建議用於線上展示)"
echo "  3) PPTX (建議用於進一步編輯)"
echo "  4) 全部格式"
echo ""
read -p "請輸入選項 (1-4): " choice

case $choice in
    1)
        echo ""
        echo "📄 轉換為 PDF..."
        marp "$INPUT_FILE" -o "RFQ簡報.pdf" --allow-local-files
        echo "✅ 完成: RFQ簡報.pdf"
        open "RFQ簡報.pdf" 2>/dev/null || echo "請手動開啟 RFQ簡報.pdf"
        ;;
    2)
        echo ""
        echo "🌐 轉換為 HTML..."
        marp "$INPUT_FILE" -o "RFQ簡報.html" --allow-local-files
        echo "✅ 完成: RFQ簡報.html"
        open "RFQ簡報.html" 2>/dev/null || echo "請手動開啟 RFQ簡報.html"
        ;;
    3)
        echo ""
        echo "📊 轉換為 PPTX..."
        marp "$INPUT_FILE" -o "RFQ簡報.pptx" --allow-local-files
        echo "✅ 完成: RFQ簡報.pptx"
        open "RFQ簡報.pptx" 2>/dev/null || echo "請手動開啟 RFQ簡報.pptx"
        ;;
    4)
        echo ""
        echo "📄 轉換為 PDF..."
        marp "$INPUT_FILE" -o "RFQ簡報.pdf" --allow-local-files
        echo "✅ PDF 完成"
        
        echo ""
        echo "🌐 轉換為 HTML..."
        marp "$INPUT_FILE" -o "RFQ簡報.html" --allow-local-files
        echo "✅ HTML 完成"
        
        echo ""
        echo "📊 轉換為 PPTX..."
        marp "$INPUT_FILE" -o "RFQ簡報.pptx" --allow-local-files
        echo "✅ PPTX 完成"
        
        echo ""
        echo "🎉 全部格式轉換完成！"
        echo ""
        echo "生成的檔案："
        echo "  - RFQ簡報.pdf"
        echo "  - RFQ簡報.html"
        echo "  - RFQ簡報.pptx"
        ;;
    *)
        echo "❌ 無效的選項"
        exit 1
        ;;
esac

echo ""
echo "🎉 轉換完成！"

