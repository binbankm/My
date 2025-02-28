// 在文件顶部添加版本信息后台密码（不可为空）
const VERSION = "1.6.8";

// 自定义标题
const CUSTOM_TITLE = "域名管理";

// 在这里设置你的 Cloudflare API Token
const CF_API_KEY = "A0yVNlSDEPe4vjt6FbfKPAAaTd11sO_ZIxdEuZqH";

// 自建 WHOIS 代理服务地址
const WHOIS_PROXY_URL = "https://whois.lbyan.us.kg";

// 访问密码（可为空）
const ACCESS_PASSWORD = "lbyan";

// 后台密码（不可为空）
const ADMIN_PASSWORD = "a4168521";

// KV 命名空间绑定名称
const KV_NAMESPACE = DOMAIN_INFO;

// 增强样式
const enhancedStyles = `
<style>
:root {
  --primary-color: #2c3e50;
  --secondary-color: #3498db;
  --success-color: #27ae60;
  --warning-color: #f1c40f;
  --danger-color: #e74c3c;
  --text-color: #34495e;
  --border-radius: 8px;
  --box-shadow: 0 2px 15px rgba(0,0,0,0.1);
}

body {
  font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
  line-height: 1.6;
  color: var(--text-color);
  background: #f8f9fa;
  margin: 0;
  padding: 20px;
}

.container {
  max-width: 1800px;
  margin: 0 auto;
  padding: 2rem 1rem;
  background: white;
  box-shadow: var(--box-shadow);
  border-radius: var(--border-radius);
  position: relative;
  min-height: calc(100vh - 120px);
}

.table-wrapper {
  overflow-x: auto;
  border-radius: var(--border-radius);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  margin: 2rem 0;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: white;
}

th {
  background: var(--primary-color);
  color: white;
  font-weight: 600;
  padding: 1rem;
  position: sticky;
  top: 0;
  z-index: 2;
}

td {
  padding: 1rem;
  border-bottom: 1px solid #eee;
  transition: background 0.2s;
}

tr:hover td {
  background: #f8f9fa;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.progress-bar {
  height: 8px;
  background: #eee;
  border-radius: 4px;
  overflow: hidden;
}

.progress {
  height: 100%;
  transition: width 0.5s ease;
}

button {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
}

.edit-btn {
  background: var(--secondary-color);
  color: white;
}

.delete-btn {
  background: var(--danger-color);
  color: white;
}

.site-footer {
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  background-color: #f8f9fa;
  color: #6c757d;
  text-align: center;
  padding: 10px 0;
  font-size: 14px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.footer-separator {
  margin: 0 10px;
}

@media (max-width: 768px) {
  .container {
    padding: 1rem;
    border-radius: 0;
  }
  
  th, td {
    padding: 0.75rem;
    font-size: 0.9em;
  }

  .mobile-hidden {
    display: none;
  }

  button {
    padding: 0.5rem;
    font-size: 0.85em;
  }
}

/* 加载动画 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255,255,255,0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.loading-spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 确认对话框 */
.confirmation-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog-content {
  background: white;
  padding: 2rem;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  max-width: 400px;
  width: 90%;
}

.dialog-buttons {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

/* Toast通知 */
.toast {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--primary-color);
  color: white;
  padding: 1rem 2rem;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  opacity: 0;
  transition: opacity 0.3s;
}

.toast.visible {
  opacity: 1;
}
</style>
`;

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
});

// ... [保留原有的 handleRequest 和路由处理函数] ...

function generateLoginHTML(title, action, errorMessage = "") {
  return `
  <!DOCTYPE html>
  <html lang="zh-CN">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title} - ${CUSTOM_TITLE}</title>
    ${enhancedStyles}
  </head>
  <body>
    <div class="container">
      <div class="login-box">
        <h1>${title}</h1>
        ${errorMessage ? `<div class="error-message">${errorMessage}</div>` : ''}
        <form method="POST" action="${action}">
          <input type="password" name="password" placeholder="请输入密码" required>
          <button type="submit">登录</button>
        </form>
      </div>
    </div>
  </body>
  </html>
  `;
}

function generateHTML(domains, isAdmin) {
  const categorizedDomains = categorizeDomains(domains);
  
  const generateDomainRows = (domainList) => domainList.map(info => {
    const today = new Date();
    const expirationDate = new Date(info.expirationDate);
    const daysRemaining = info.expirationDate === 'Unknown' ? 'N/A' : Math.ceil((expirationDate - today) / (1000 * 60 * 60 * 24));
    const totalDays = info.registrationDate === 'Unknown' || info.expirationDate === 'Unknown' ? 'N/A' : Math.ceil((expirationDate - new Date(info.registrationDate)) / (1000 * 60 * 60 * 24));
    const progressPercentage = isNaN(daysRemaining) || isNaN(totalDays) ? 0 : 100 - (daysRemaining / totalDays * 100);
    
    return `
      <tr>
        <td><span class="status-dot" style="background:${getStatusColor(daysRemaining)}"></span></td>
        <td>${info.domain}</td>
        <td class="mobile-hidden">${info.system}</td>
        <td class="mobile-hidden">${info.registrar}</td>
        <td>${info.registrationDate}</td>
        <td>${info.expirationDate}</td>
        <td>${daysRemaining}</td>
        <td>
          <div class="progress-bar">
            <div class="progress" style="width:${progressPercentage}%"></div>
          </div>
        </td>
        ${isAdmin ? `
        <td>
          <button class="edit-btn" onclick="editDomain('${info.domain}', this)">编辑</button>
          <button class="delete-btn" onclick="deleteDomain('${info.domain}')">删除</button>
        </td>` : ''}
      </tr>
    `;
  }).join('');

  return `
  <!DOCTYPE html>
  <html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${CUSTOM_TITLE}${isAdmin ? ' - 后台管理' : ''}</title>
    ${enhancedStyles}
    <script>
    let isLoading = false;

    async function showLoading() {
      isLoading = true;
      const loader = document.createElement('div');
      loader.className = 'loading-overlay';
      loader.innerHTML = '<div class="loading-spinner"></div>';
      document.body.appendChild(loader);
    }

    async function hideLoading() {
      isLoading = false;
      document.querySelector('.loading-overlay')?.remove();
    }

    async function confirmAction(message) {
      return new Promise(resolve => {
        const dialog = document.createElement('div');
        dialog.className = 'confirmation-dialog';
        dialog.innerHTML = \`
          <div class="dialog-content">
            <p>${message}</p>
            <div class="dialog-buttons">
              <button class="confirm-btn">确定</button>
              <button class="cancel-btn">取消</button>
            </div>
          </div>
        \`;
        document.body.appendChild(dialog);

        dialog.querySelector('.confirm-btn').onclick = () => {
          dialog.remove();
          resolve(true);
        };
        dialog.querySelector('.cancel-btn').onclick = () => {
          dialog.remove();
          resolve(false);
        };
      });
    }

    async function editDomain(domain, button) {
      if (isLoading) return;
      
      const row = button.closest('tr');
      const cells = row.querySelectorAll('td:not(:last-child)');
      
      if (button.textContent === '编辑') {
        button.textContent = '保存';
        button.classList.add('saving');
        cells.forEach((cell, index) => {
          if(index > 0 && index < 5) { // 只允许编辑特定列
            const input = document.createElement('input');
            input.value = cell.textContent.trim();
            cell.innerHTML = '';
            cell.appendChild(input);
          }
        });
      } else {
        const confirmed = await confirmAction('确定要保存修改吗？');
        if (!confirmed) return;
        
        await showLoading();
        try {
          const updatedData = {
            domain: domain,
            registrar: cells[2].querySelector('input').value,
            registrationDate: cells[3].querySelector('input').value,
            expirationDate: cells[4].querySelector('input').value
          };

          const response = await fetch('/api/update', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': 'Basic ' + btoa(':${ADMIN_PASSWORD}')
            },
            body: JSON.stringify(updatedData)
          });

          if (response.ok) {
            cells.forEach(cell => {
              if (cell.querySelector('input')) {
                cell.textContent = cell.querySelector('input').value;
              }
            });
            showToast('更新成功', 'success');
          } else {
            throw new Error('更新失败');
          }
        } catch (error) {
          showToast('更新失败: ' + error.message, 'error');
          row.classList.add('error-shake');
          setTimeout(() => row.classList.remove('error-shake'), 500);
        } finally {
          button.textContent = '编辑';
          button.classList.remove('saving');
          await hideLoading();
        }
      }
    }

    async function deleteDomain(domain) {
      const confirmed = await confirmAction('确定要删除这个域名吗？');
      if (!confirmed) return;
      
      await showLoading();
      try {
        const response = await fetch('/api/update', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Basic ' + btoa(':${ADMIN_PASSWORD}')
          },
          body: JSON.stringify({ action: 'delete', domain: domain })
        });

        if (response.ok) {
          document.querySelector(\`tr[data-domain="${domain}"]\`)?.remove();
          showToast('删除成功', 'success');
        } else {
          throw new Error('删除失败');
        }
      } catch (error) {
        showToast('删除失败: ' + error.message, 'error');
      } finally {
        await hideLoading();
      }
    }

    function showToast(message, type) {
      const toast = document.createElement('div');
      toast.className = \`toast \${type}\`;
      toast.textContent = message;
      document.body.appendChild(toast);
      setTimeout(() => {
        toast.classList.add('visible');
        setTimeout(() => {
          toast.remove();
        }, 3000);
      }, 100);
    }
    </script>
  </head>
  <body>
    <div class="container">
      <header class="page-header">
        <h1>${CUSTOM_TITLE}${isAdmin ? ' - 后台管理' : ''}</h1>
        <div class="admin-controls">
          ${isAdmin ? 
            '<a href="/" class="admin-link">返回前台</a>' : 
            '<a href="/admin" class="admin-link">后台管理</a>'
          }
        </div>
      </header>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th class="status-column">状态</th>
              <th>域名</th>
              <th class="mobile-hidden">系统</th>
              <th class="mobile-hidden">注册商</th>
              <th>注册日期</th>
              <th>到期日期</th>
              <th>剩余天数</th>
              <th>进度</th>
              ${isAdmin ? '<th>操作</th>' : ''}
            </tr>
          </thead>
          <tbody>
            ${generateDomainRows([...categorizedDomains.cfTopLevel, ...categorizedDomains.cfSecondLevelAndCustom])}
          </tbody>
        </table>
      </div>
    </div>
    ${footerHTML}
  </body>
  </html>
  `;
}

// ... [保留其他原有函数] ...
