
import type { DefineComponent, SlotsType } from 'vue'
type IslandComponent<T> = DefineComponent<{}, {refresh: () => Promise<void>}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, SlotsType<{ fallback: { error: unknown } }>> & T

type HydrationStrategies = {
  hydrateOnVisible?: IntersectionObserverInit | true
  hydrateOnIdle?: number | true
  hydrateOnInteraction?: keyof HTMLElementEventMap | Array<keyof HTMLElementEventMap> | true
  hydrateOnMediaQuery?: string
  hydrateAfter?: number
  hydrateWhen?: boolean
  hydrateNever?: true
}
type LazyComponent<T> = DefineComponent<HydrationStrategies, {}, {}, {}, {}, {}, {}, { hydrated: () => void }> & T

interface _GlobalComponents {
  AuthModal: typeof import("../../components/auth/AuthModal.vue")['default']
  AuthCaptchaModal: typeof import("../../components/auth/CaptchaModal.vue")['default']
  AuthLogin: typeof import("../../components/auth/Login.vue")['default']
  AuthPassword: typeof import("../../components/auth/Password.vue")['default']
  AuthRegister: typeof import("../../components/auth/Register.vue")['default']
  AuthResetPassword: typeof import("../../components/auth/ResetPassword.vue")['default']
  CodeEditor: typeof import("../../components/code/CodeEditor.vue")['default']
  CodeSidebar: typeof import("../../components/code/CodeSidebar.vue")['default']
  CodeConfigPanel: typeof import("../../components/code/ConfigPanel.vue")['default']
  CodeFlowSteps: typeof import("../../components/code/FlowSteps.vue")['default']
  CodeGeneratingOverlay: typeof import("../../components/code/GeneratingOverlay.vue")['default']
  CodeMonacoEditorCore: typeof import("../../components/code/MonacoEditorCore.vue")['default']
  CommonFlutterTemplate: typeof import("../../components/common/FlutterTemplate.vue")['default']
  DialogModel: typeof import("../../components/dialog/DialogModel.vue")['default']
  LandingDiagnosticBubbles: typeof import("../../components/landing/DiagnosticBubbles.vue")['default']
  LandingHeroSection: typeof import("../../components/landing/HeroSection.vue")['default']
  LandingSSEGenerator: typeof import("../../components/landing/SSEGenerator.vue")['default']
  LayoutAppFooter: typeof import("../../components/layout/AppFooter.vue")['default']
  LayoutAppHeader: typeof import("../../components/layout/AppHeader.vue")['default']
  LayoutAppSidebar: typeof import("../../components/layout/AppSidebar.vue")['default']
  LayoutParticlesBackground: typeof import("../../components/layout/ParticlesBackground.vue")['default']
  PreviewTempLateComputer: typeof import("../../components/previewTempLate/computer.vue")['default']
  PreviewTempLate: typeof import("../../components/previewTempLate/index.vue")['default']
  PreviewTempLateMobile: typeof import("../../components/previewTempLate/mobile.vue")['default']
  TasksFilter: typeof import("../../components/tasks/TasksFilter.vue")['default']
  TasksSidebar: typeof import("../../components/tasks/TasksSidebar.vue")['default']
  TasksStatusOverview: typeof import("../../components/tasks/TasksStatusOverview.vue")['default']
  TasksTable: typeof import("../../components/tasks/TasksTable.vue")['default']
  UploadDescEditorModal: typeof import("../../components/upload/DescEditorModal.vue")['default']
  UploadZone: typeof import("../../components/upload/UploadZone.vue")['default']
  NuxtWelcome: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/welcome.vue")['default']
  NuxtLayout: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-layout")['default']
  NuxtErrorBoundary: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
  ClientOnly: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/client-only")['default']
  DevOnly: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/dev-only")['default']
  ServerPlaceholder: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/server-placeholder")['default']
  NuxtLink: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-link")['default']
  NuxtLoadingIndicator: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
  NuxtTime: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
  NuxtRouteAnnouncer: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
  NuxtAnnouncer: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-announcer")['default']
  NuxtImg: typeof import("../../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtImg.vue")['default']
  NuxtPicture: typeof import("../../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtPicture.vue")['default']
  ColorScheme: typeof import("../../node_modules/.pnpm/@nuxtjs+color-mode@4.0.1_magicast@0.5.3/node_modules/@nuxtjs/color-mode/dist/runtime/component.vue")['default']
  NuxtPage: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/pages/runtime/page")['default']
  NoScript: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['NoScript']
  Link: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Link']
  Base: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Base']
  Title: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Title']
  Meta: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Meta']
  Style: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Style']
  Head: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Head']
  Html: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Html']
  Body: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Body']
  NuxtIsland: typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-island")['default']
  LazyAuthModal: LazyComponent<typeof import("../../components/auth/AuthModal.vue")['default']>
  LazyAuthCaptchaModal: LazyComponent<typeof import("../../components/auth/CaptchaModal.vue")['default']>
  LazyAuthLogin: LazyComponent<typeof import("../../components/auth/Login.vue")['default']>
  LazyAuthPassword: LazyComponent<typeof import("../../components/auth/Password.vue")['default']>
  LazyAuthRegister: LazyComponent<typeof import("../../components/auth/Register.vue")['default']>
  LazyAuthResetPassword: LazyComponent<typeof import("../../components/auth/ResetPassword.vue")['default']>
  LazyCodeEditor: LazyComponent<typeof import("../../components/code/CodeEditor.vue")['default']>
  LazyCodeSidebar: LazyComponent<typeof import("../../components/code/CodeSidebar.vue")['default']>
  LazyCodeConfigPanel: LazyComponent<typeof import("../../components/code/ConfigPanel.vue")['default']>
  LazyCodeFlowSteps: LazyComponent<typeof import("../../components/code/FlowSteps.vue")['default']>
  LazyCodeGeneratingOverlay: LazyComponent<typeof import("../../components/code/GeneratingOverlay.vue")['default']>
  LazyCodeMonacoEditorCore: LazyComponent<typeof import("../../components/code/MonacoEditorCore.vue")['default']>
  LazyCommonFlutterTemplate: LazyComponent<typeof import("../../components/common/FlutterTemplate.vue")['default']>
  LazyDialogModel: LazyComponent<typeof import("../../components/dialog/DialogModel.vue")['default']>
  LazyLandingDiagnosticBubbles: LazyComponent<typeof import("../../components/landing/DiagnosticBubbles.vue")['default']>
  LazyLandingHeroSection: LazyComponent<typeof import("../../components/landing/HeroSection.vue")['default']>
  LazyLandingSSEGenerator: LazyComponent<typeof import("../../components/landing/SSEGenerator.vue")['default']>
  LazyLayoutAppFooter: LazyComponent<typeof import("../../components/layout/AppFooter.vue")['default']>
  LazyLayoutAppHeader: LazyComponent<typeof import("../../components/layout/AppHeader.vue")['default']>
  LazyLayoutAppSidebar: LazyComponent<typeof import("../../components/layout/AppSidebar.vue")['default']>
  LazyLayoutParticlesBackground: LazyComponent<typeof import("../../components/layout/ParticlesBackground.vue")['default']>
  LazyPreviewTempLateComputer: LazyComponent<typeof import("../../components/previewTempLate/computer.vue")['default']>
  LazyPreviewTempLate: LazyComponent<typeof import("../../components/previewTempLate/index.vue")['default']>
  LazyPreviewTempLateMobile: LazyComponent<typeof import("../../components/previewTempLate/mobile.vue")['default']>
  LazyTasksFilter: LazyComponent<typeof import("../../components/tasks/TasksFilter.vue")['default']>
  LazyTasksSidebar: LazyComponent<typeof import("../../components/tasks/TasksSidebar.vue")['default']>
  LazyTasksStatusOverview: LazyComponent<typeof import("../../components/tasks/TasksStatusOverview.vue")['default']>
  LazyTasksTable: LazyComponent<typeof import("../../components/tasks/TasksTable.vue")['default']>
  LazyUploadDescEditorModal: LazyComponent<typeof import("../../components/upload/DescEditorModal.vue")['default']>
  LazyUploadZone: LazyComponent<typeof import("../../components/upload/UploadZone.vue")['default']>
  LazyNuxtWelcome: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/welcome.vue")['default']>
  LazyNuxtLayout: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
  LazyNuxtErrorBoundary: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
  LazyClientOnly: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/client-only")['default']>
  LazyDevOnly: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/dev-only")['default']>
  LazyServerPlaceholder: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/server-placeholder")['default']>
  LazyNuxtLink: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-link")['default']>
  LazyNuxtLoadingIndicator: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
  LazyNuxtTime: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
  LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
  LazyNuxtAnnouncer: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-announcer")['default']>
  LazyNuxtImg: LazyComponent<typeof import("../../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtImg.vue")['default']>
  LazyNuxtPicture: LazyComponent<typeof import("../../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtPicture.vue")['default']>
  LazyColorScheme: LazyComponent<typeof import("../../node_modules/.pnpm/@nuxtjs+color-mode@4.0.1_magicast@0.5.3/node_modules/@nuxtjs/color-mode/dist/runtime/component.vue")['default']>
  LazyNuxtPage: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/pages/runtime/page")['default']>
  LazyNoScript: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['NoScript']>
  LazyLink: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Link']>
  LazyBase: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Base']>
  LazyTitle: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Title']>
  LazyMeta: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Meta']>
  LazyStyle: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Style']>
  LazyHead: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Head']>
  LazyHtml: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Html']>
  LazyBody: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Body']>
  LazyNuxtIsland: LazyComponent<typeof import("../../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-island")['default']>
}

declare module 'vue' {
  export interface GlobalComponents extends _GlobalComponents { }
}

export {}
