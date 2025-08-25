function generateEmail() {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    const domains = ['gmail.com', '163.com', '126.com', 'qq.com', 'outlook.com', 'hotmail.com'];
    const name = Array.from({ length: 10 }, () => chars[Math.floor(Math.random() * chars.length)]).join('');
    const domain = domains[Math.floor(Math.random() * domains.length)];
    return `${name}@${domain}`;
}

function generatePassword(length = 10) {
    const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    return Array.from({ length }, () => chars[Math.floor(Math.random() * chars.length)]).join('');
}

export default {
    async fetch(request, env) {
        // 修改为目标机场的注册和订阅基地址
        const baseUrl = "https://your-airport-domain.com"; // 注册 API 域名
        const subscribeBase = "https://subscribe-domain.com"; // 订阅 API 域名
        
        const apiBase = `${baseUrl}/api/v1`;
        const userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36';
        const headersBase = {
            'Content-Type': 'application/json',
            'Origin': baseUrl,
            'User-Agent': userAgent
        };
        
        const email = generateEmail();
        const password = generatePassword();
        
        try {
            
            const registerData = JSON.stringify({ email, password });
            const registerHeaders = { ...headersBase, 'Referer': `${baseUrl}/#/register` };
            const registerResponse = await fetch(`${apiBase}/passport/auth/register`, {
                method: 'POST',
                headers: registerHeaders,
                body: registerData
            });
            if (!registerResponse.ok) {
                return new Response(`注册失败: ${await registerResponse.text()}`, { status: 500 });
            }
            const registerResult = await registerResponse.json();
            const subToken = registerResult.data?.token || registerResult.data?.auth_data?.token;
            if (!subToken) {
                return new Response(`注册无 token: ${JSON.stringify(registerResult)}`, { status: 500 });
            }
            
            const finalUrl = `${subscribeBase}/api/v1/client/subscribe?token=${subToken}`;
            const contentResponse = await fetch(finalUrl, {
                method: 'GET',
                headers: { 'User-Agent': userAgent }
            });
            if (!contentResponse.ok) {
                return new Response(`获取订阅内容失败: ${await contentResponse.text()}`, { status: 500 });
            }
            const content = await contentResponse.text();
            
            try {
                const decodedContent = atob(content);
                return new Response(decodedContent, {
                    headers: { 'Content-Type': 'text/plain; charset=utf-8' }
                });
            } catch {
                return new Response(content, {
                    headers: { 'Content-Type': 'text/plain; charset=utf-8' }
                });
            }
        } catch (e) {
            return new Response(`整体错误: ${e.message}`, { status: 500 });
        }
    }
};