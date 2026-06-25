
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


export const AuthModal: typeof import("../components/auth/AuthModal.vue")['default']
export const AuthCaptchaModal: typeof import("../components/auth/CaptchaModal.vue")['default']
export const AuthLogin: typeof import("../components/auth/Login.vue")['default']
export const AuthPassword: typeof import("../components/auth/Password.vue")['default']
export const AuthRegister: typeof import("../components/auth/Register.vue")['default']
export const CodeEditor: typeof import("../components/code/CodeEditor.vue")['default']
export const CodeSidebar: typeof import("../components/code/CodeSidebar.vue")['default']
export const CodeConfigPanel: typeof import("../components/code/ConfigPanel.vue")['default']
export const CodeFlowSteps: typeof import("../components/code/FlowSteps.vue")['default']
export const CodeGeneratingOverlay: typeof import("../components/code/GeneratingOverlay.vue")['default']
export const CodeMonacoEditorCore: typeof import("../components/code/MonacoEditorCore.vue")['default']
export const CommonFlutterTemplate: typeof import("../components/common/FlutterTemplate.vue")['default']
export const DialogModel: typeof import("../components/dialog/DialogModel.vue")['default']
export const LandingDiagnosticBubbles: typeof import("../components/landing/DiagnosticBubbles.vue")['default']
export const LandingHeroSection: typeof import("../components/landing/HeroSection.vue")['default']
export const LandingSSEGenerator: typeof import("../components/landing/SSEGenerator.vue")['default']
export const LayoutAppFooter: typeof import("../components/layout/AppFooter.vue")['default']
export const LayoutAppHeader: typeof import("../components/layout/AppHeader.vue")['default']
export const LayoutAppSidebar: typeof import("../components/layout/AppSidebar.vue")['default']
export const LayoutParticlesBackground: typeof import("../components/layout/ParticlesBackground.vue")['default']
export const PreviewTempLateComputer: typeof import("../components/previewTempLate/computer.vue")['default']
export const PreviewTempLate: typeof import("../components/previewTempLate/index.vue")['default']
export const PreviewTempLateMobile: typeof import("../components/previewTempLate/mobile.vue")['default']
export const TasksSidebar: typeof import("../components/tasks/TasksSidebar.vue")['default']
export const UploadDescEditorModal: typeof import("../components/upload/DescEditorModal.vue")['default']
export const UploadZone: typeof import("../components/upload/UploadZone.vue")['default']
export const NuxtWelcome: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/welcome.vue")['default']
export const NuxtLayout: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-layout")['default']
export const NuxtErrorBoundary: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
export const ClientOnly: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/client-only")['default']
export const DevOnly: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/dev-only")['default']
export const ServerPlaceholder: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/server-placeholder")['default']
export const NuxtLink: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-link")['default']
export const NuxtLoadingIndicator: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
export const NuxtTime: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
export const NuxtRouteAnnouncer: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
export const NuxtAnnouncer: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-announcer")['default']
export const NuxtImg: typeof import("../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtImg.vue")['default']
export const NuxtPicture: typeof import("../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtPicture.vue")['default']
export const ColorScheme: typeof import("../node_modules/.pnpm/@nuxtjs+color-mode@4.0.1_magicast@0.5.3/node_modules/@nuxtjs/color-mode/dist/runtime/component.vue")['default']
export const NuxtPage: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/pages/runtime/page")['default']
export const NoScript: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['NoScript']
export const Link: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Link']
export const Base: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Base']
export const Title: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Title']
export const Meta: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Meta']
export const Style: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Style']
export const Head: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Head']
export const Html: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Html']
export const Body: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Body']
export const NuxtIsland: typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-island")['default']
export const LazyAuthModal: LazyComponent<typeof import("../components/auth/AuthModal.vue")['default']>
export const LazyAuthCaptchaModal: LazyComponent<typeof import("../components/auth/CaptchaModal.vue")['default']>
export const LazyAuthLogin: LazyComponent<typeof import("../components/auth/Login.vue")['default']>
export const LazyAuthPassword: LazyComponent<typeof import("../components/auth/Password.vue")['default']>
export const LazyAuthRegister: LazyComponent<typeof import("../components/auth/Register.vue")['default']>
export const LazyCodeEditor: LazyComponent<typeof import("../components/code/CodeEditor.vue")['default']>
export const LazyCodeSidebar: LazyComponent<typeof import("../components/code/CodeSidebar.vue")['default']>
export const LazyCodeConfigPanel: LazyComponent<typeof import("../components/code/ConfigPanel.vue")['default']>
export const LazyCodeFlowSteps: LazyComponent<typeof import("../components/code/FlowSteps.vue")['default']>
export const LazyCodeGeneratingOverlay: LazyComponent<typeof import("../components/code/GeneratingOverlay.vue")['default']>
export const LazyCodeMonacoEditorCore: LazyComponent<typeof import("../components/code/MonacoEditorCore.vue")['default']>
export const LazyCommonFlutterTemplate: LazyComponent<typeof import("../components/common/FlutterTemplate.vue")['default']>
export const LazyDialogModel: LazyComponent<typeof import("../components/dialog/DialogModel.vue")['default']>
export const LazyLandingDiagnosticBubbles: LazyComponent<typeof import("../components/landing/DiagnosticBubbles.vue")['default']>
export const LazyLandingHeroSection: LazyComponent<typeof import("../components/landing/HeroSection.vue")['default']>
export const LazyLandingSSEGenerator: LazyComponent<typeof import("../components/landing/SSEGenerator.vue")['default']>
export const LazyLayoutAppFooter: LazyComponent<typeof import("../components/layout/AppFooter.vue")['default']>
export const LazyLayoutAppHeader: LazyComponent<typeof import("../components/layout/AppHeader.vue")['default']>
export const LazyLayoutAppSidebar: LazyComponent<typeof import("../components/layout/AppSidebar.vue")['default']>
export const LazyLayoutParticlesBackground: LazyComponent<typeof import("../components/layout/ParticlesBackground.vue")['default']>
export const LazyPreviewTempLateComputer: LazyComponent<typeof import("../components/previewTempLate/computer.vue")['default']>
export const LazyPreviewTempLate: LazyComponent<typeof import("../components/previewTempLate/index.vue")['default']>
export const LazyPreviewTempLateMobile: LazyComponent<typeof import("../components/previewTempLate/mobile.vue")['default']>
export const LazyTasksSidebar: LazyComponent<typeof import("../components/tasks/TasksSidebar.vue")['default']>
export const LazyUploadDescEditorModal: LazyComponent<typeof import("../components/upload/DescEditorModal.vue")['default']>
export const LazyUploadZone: LazyComponent<typeof import("../components/upload/UploadZone.vue")['default']>
export const LazyNuxtWelcome: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/welcome.vue")['default']>
export const LazyNuxtLayout: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
export const LazyNuxtErrorBoundary: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
export const LazyClientOnly: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/client-only")['default']>
export const LazyDevOnly: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/dev-only")['default']>
export const LazyServerPlaceholder: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/server-placeholder")['default']>
export const LazyNuxtLink: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-link")['default']>
export const LazyNuxtLoadingIndicator: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
export const LazyNuxtTime: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
export const LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
export const LazyNuxtAnnouncer: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-announcer")['default']>
export const LazyNuxtImg: LazyComponent<typeof import("../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtImg.vue")['default']>
export const LazyNuxtPicture: LazyComponent<typeof import("../node_modules/.pnpm/@nuxt+image@2.0.0_db0@0.3.4_98b0fc79736c87c2a668df3f5ed3e552/node_modules/@nuxt/image/dist/runtime/components/NuxtPicture.vue")['default']>
export const LazyColorScheme: LazyComponent<typeof import("../node_modules/.pnpm/@nuxtjs+color-mode@4.0.1_magicast@0.5.3/node_modules/@nuxtjs/color-mode/dist/runtime/component.vue")['default']>
export const LazyNuxtPage: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/pages/runtime/page")['default']>
export const LazyNoScript: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['NoScript']>
export const LazyLink: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Link']>
export const LazyBase: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Base']>
export const LazyTitle: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Title']>
export const LazyMeta: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Meta']>
export const LazyStyle: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Style']>
export const LazyHead: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Head']>
export const LazyHtml: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Html']>
export const LazyBody: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/head/runtime/components")['Body']>
export const LazyNuxtIsland: LazyComponent<typeof import("../node_modules/.pnpm/nuxt@4.4.8_@babel+plugin-sy_99ff0ff52633701d5303af6e8aaecd08/node_modules/nuxt/dist/app/components/nuxt-island")['default']>

export const componentNames: string[]
