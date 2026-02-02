/**
 * 权限刷新脚本
 *
 * 使用方法：
 * 1. 打开浏览器，按 F12 打开开发者工具
 * 2. 切换到 Console 标签
 * 3. 复制下面的脚本并粘贴到控制台
 * 4. 按回车执行
 */

(function() {
  console.log('=== 权限刷新工具 ===\n');

  // 1. 显示当前权限
  console.log('📋 当前本地存储的权限:');
  const currentPermissions = JSON.parse(localStorage.getItem('permissions') || '[]');
  console.log('权限数量:', currentPermissions.length);
  console.log('物资相关权限:', currentPermissions.filter(p => p.startsWith('material_')));

  // 2. 显示当前用户信息
  console.log('\n👤 当前用户信息:');
  const currentUser = JSON.parse(localStorage.getItem('user') || 'null');
  console.log('用户名:', currentUser?.username);
  console.log('角色:', currentUser?.roles?.map(r => r.name).join(', '));

  // 3. 检查物资管理权限
  console.log('\n🔍 权限检查结果:');
  const requiredPerms = ['material_view', 'material_create', 'material_edit', 'material_delete', 'material_import', 'material_export', 'material_in'];
  const hasMaterialPermission = requiredPerms.some(p => currentPermissions.includes(p));

  if (hasMaterialPermission) {
    console.log('✅ 您拥有物资管理权限');
    requiredPerms.forEach(p => {
      const has = currentPermissions.includes(p);
      console.log(`  ${has ? '✅' : '❌'} ${p}`);
    });
  } else {
    console.log('❌ 您没有物资管理权限');
    console.log('\n💡 解决方案:');
    console.log('1. 点击右上角用户头像');
    console.log('2. 选择"退出登录"');
    console.log('3. 重新登录');
    console.log('4. 刷新页面');
  }

  // 4. 提供一键清除权限的选项
  console.log('\n🛠️  快速操作:');
  console.log('如需强制刷新权限，请执行以下命令:');
  console.log('localStorage.clear(); location.href="/login"');

  console.log('\n===================\n');
})();
