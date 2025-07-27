import { createRouter, createWebHistory } from 'vue-router'

// 首页
const Home = () => import('@/views/Home.vue')

// 静态页面
const EmailService = () => import('@/views/EmailService.vue')
const ApiDocs = () => import('@/views/ApiDocs.vue')
const HelpCenter = () => import('@/views/HelpCenter.vue')
const Contact = () => import('@/views/Contact.vue')
const TechSupport = () => import('@/views/TechSupport.vue')
const PrivacyPolicy = () => import('@/views/PrivacyPolicy.vue')
const TermsOfService = () => import('@/views/TermsOfService.vue')
const CookiePolicy = () => import('@/views/CookiePolicy.vue')

// 用户端页面
const UserLogin = () => import('@/views/user/Login.vue')
const UserRegister = () => import('@/views/user/Register.vue')
const UserDashboard = () => import('@/views/user/Dashboard.vue')
const UserInbox = () => import('@/views/user/Inbox.vue')
const UserOutbox = () => import('@/views/user/Outbox.vue')
const UserCompose = () => import('@/views/user/Compose.vue')
const UserSettings = () => import('@/views/user/Settings.vue')
const UserMailboxes = () => import('@/views/user/Mailboxes.vue')
const UserDomains = () => import('@/views/user/Domains.vue')
const UserForwardRules = () => import('@/views/user/ForwardRules.vue')

// 管理员页面
const AdminLogin = () => import('@/views/admin/Login.vue')
const AdminDashboard = () => import('@/views/admin/Dashboard.vue')
const AdminUsers = () => import('@/views/admin/Users.vue')
const AdminMailboxes = () => import('@/views/admin/Mailboxes.vue')
const AdminDomains = () => import('@/views/admin/Domains.vue')
const AdminCaptcha = () => import('@/views/admin/Captcha.vue')

// 布局组件
const UserLayout = () => import('@/layouts/UserLayout.vue')
const AdminLayout = () => import('@/layouts/AdminLayout.vue')

const routes = [
  // 首页
  {
    path: '/',
    name: 'Home',
    component: Home
  },

  // 静态页面路由
  {
    path: '/email-service',
    name: 'EmailService',
    component: EmailService
  },
  {
    path: '/api-docs',
    name: 'ApiDocs',
    component: ApiDocs
  },
  {
    path: '/help-center',
    name: 'HelpCenter',
    component: HelpCenter
  },
  {
    path: '/contact',
    name: 'Contact',
    component: Contact
  },
  {
    path: '/tech-support',
    name: 'TechSupport',
    component: TechSupport
  },
  {
    path: '/privacy-policy',
    name: 'PrivacyPolicy',
    component: PrivacyPolicy
  },
  {
    path: '/terms-of-service',
    name: 'TermsOfService',
    component: TermsOfService
  },
  {
    path: '/cookie-policy',
    name: 'CookiePolicy',
    component: CookiePolicy
  },

  // 用户端路由
  {
    path: '/user',
    children: [
      {
        path: 'login',
        name: 'UserLogin',
        component: UserLogin
      },
      {
        path: 'register',
        name: 'UserRegister',
        component: UserRegister
      },
      {
        path: '',
        component: UserLayout,
        meta: { requiresAuth: true },
        children: [
          {
            path: 'dashboard',
            name: 'UserDashboard',
            component: UserDashboard
          },
          {
            path: 'inbox',
            name: 'UserInbox',
            component: UserInbox
          },
          {
            path: 'outbox',
            name: 'UserOutbox',
            component: UserOutbox
          },
          {
            path: 'compose',
            name: 'UserCompose',
            component: UserCompose
          },
          {
            path: 'forward-rules',
            name: 'UserForwardRules',
            component: UserForwardRules
          },
          {
            path: 'mailboxes',
            name: 'UserMailboxes',
            component: UserMailboxes
          },
          {
            path: 'domains',
            name: 'UserDomains',
            component: UserDomains
          },
          {
            path: 'settings',
            name: 'UserSettings',
            component: UserSettings
          }
        ]
      }
    ]
  },

  // 管理员路由
  {
    path: '/admin',
    children: [
      {
        path: 'login',
        name: 'AdminLogin',
        component: AdminLogin
      },
      {
        path: '',
        component: AdminLayout,
        meta: { requiresAdminAuth: true },
        children: [
          {
            path: 'dashboard',
            name: 'AdminDashboard',
            component: AdminDashboard
          },
          {
            path: 'users',
            name: 'AdminUsers',
            component: AdminUsers
          },
          {
            path: 'mailboxes',
            name: 'AdminMailboxes',
            component: AdminMailboxes
          },
          {
            path: 'domains',
            name: 'AdminDomains',
            component: AdminDomains
          },
          {
            path: 'captcha',
            name: 'AdminCaptcha',
            component: AdminCaptcha
          },
          {
            path: 'api-test',
            name: 'AdminApiTest',
            component: () => import('@/views/admin/ApiTest.vue')
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userToken = localStorage.getItem('user_token')
  const adminToken = localStorage.getItem('admin_token')

  if (to.meta.requiresAuth && !userToken) {
    next('/user/login')
  } else if (to.meta.requiresAdminAuth && !adminToken) {
    next('/admin/login')
  } else {
    next()
  }
})

export default router
