#!/bin/bash

echo "🧪 测试浏览器 MCP 配置..."
echo ""

# 检查 mcpbrowser 是否安装
if command -v mcpbrowser &> /dev/null; then
    echo "✅ mcpbrowser 已安装: $(which mcpbrowser)"
    echo "   版本: $(npm list -g mcpbrowser 2>/dev/null | grep mcpbrowser)"
else
    echo "❌ mcpbrowser 未安装"
    exit 1
fi

echo ""
echo "📋 MCP 配置文件位置:"
echo "   ~/Library/Application Support/Cursor/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json"
echo ""

# 检查配置文件
CONFIG_FILE="$HOME/Library/Application Support/Cursor/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json"
if [ -f "$CONFIG_FILE" ]; then
    echo "✅ 配置文件存在"
    echo ""
    echo "📄 配置内容:"
    cat "$CONFIG_FILE" | python3 -m json.tool 2>/dev/null || cat "$CONFIG_FILE"
else
    echo "❌ 配置文件不存在"
    exit 1
fi

echo ""
echo "🔍 检查 MCP 进程..."
if ps aux | grep -i "mcpbrowser\|mcp" | grep -v grep > /dev/null; then
    echo "✅ MCP 进程正在运行"
    ps aux | grep -i "mcpbrowser\|mcp" | grep -v grep
else
    echo "⚠️  MCP 进程未运行（可能需要重启 Cursor 或手动启动）"
fi

echo ""
echo "📝 下一步："
echo "   1. 确保已重启 Cursor"
echo "   2. 在 Cursor 中按 Cmd+Shift+P，输入 'MCP' 查看相关命令"
echo "   3. 或者直接问我：'帮我访问 http://localhost:3000 并截图'"
echo ""
