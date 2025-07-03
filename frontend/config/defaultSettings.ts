import { ProLayoutProps } from '@ant-design/pro-components';

/**
 * LYSS AI Platform 默认配置
 * 包含主题、布局、颜色等基础设置
 */
const Settings: ProLayoutProps & {
  pwa?: boolean;
  logo?: string;
} = {
  navTheme: 'light',
  // LYSS AI 主色调 - 科技蓝
  colorPrimary: '#1677ff',
  layout: 'mix',
  contentWidth: 'Fluid',
  fixedHeader: false,
  fixSiderbar: true,
  colorWeak: false,
  title: 'LYSS AI Platform',
  pwa: false, // 暂时禁用PWA
  logo: '/logo.svg', // 使用本地Logo
  iconfontUrl: '',
  token: {
    // 通过token自定义样式主题
    colorBgContainer: '#ffffff',
    borderRadius: 6,
  },
};

export default Settings;
