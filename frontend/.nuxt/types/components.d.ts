
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
  AuthLogin: typeof import("../../components/auth/Login.vue")['default']
  AuthLoginModal: typeof import("../../components/auth/LoginModal.vue")['default']
  AuthRegister: typeof import("../../components/auth/Register.vue")['default']
  CodeEditor: typeof import("../../components/code/CodeEditor.vue")['default']
  CodePreview: typeof import("../../components/code/CodePreview.vue")['default']
  CodeGeneratingOverlay: typeof import("../../components/code/GeneratingOverlay.vue")['default']
  CodeMonacoEditorCore: typeof import("../../components/code/MonacoEditorCore.vue")['default']
  CodeOptimizationPanel: typeof import("../../components/code/OptimizationPanel.vue")['default']
  CodeSyntaxHighlighter: typeof import("../../components/code/SyntaxHighlighter.vue")['default']
  CommonLoadingSpinner: typeof import("../../components/common/LoadingSpinner.vue")['default']
  CommonModalDialog: typeof import("../../components/common/ModalDialog.vue")['default']
  CommonToastNotification: typeof import("../../components/common/ToastNotification.vue")['default']
  ComparisonDiffViewer: typeof import("../../components/comparison/DiffViewer.vue")['default']
  ComparisonModelComparison: typeof import("../../components/comparison/ModelComparison.vue")['default']
  ComparisonVersionTimeline: typeof import("../../components/comparison/VersionTimeline.vue")['default']
  EvaluationQualityGauge: typeof import("../../components/evaluation/QualityGauge.vue")['default']
  EvaluationScoreCard: typeof import("../../components/evaluation/ScoreCard.vue")['default']
  EvaluationSuggestionList: typeof import("../../components/evaluation/SuggestionList.vue")['default']
  LayoutAppFooter: typeof import("../../components/layout/AppFooter.vue")['default']
  LayoutAppHeader: typeof import("../../components/layout/AppHeader.vue")['default']
  LayoutCardsArea: typeof import("../../components/layout/CardsArea.vue")['default']
  LayoutHeroSection: typeof import("../../components/layout/HeroSection.vue")['default']
  OptimizationCompareModal: typeof import("../../components/optimization/CompareModal.vue")['default']
  OptimizationProgressModal: typeof import("../../components/optimization/ProgressModal.vue")['default']
  UploadDragDropZone: typeof import("../../components/upload/DragDropZone.vue")['default']
  UploadImagePreview: typeof import("../../components/upload/ImagePreview.vue")['default']
  UploadImageUploader: typeof import("../../components/upload/ImageUploader.vue")['default']
  NuxtWelcome: typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']
  NuxtLayout: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
  NuxtErrorBoundary: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
  ClientOnly: typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']
  DevOnly: typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']
  ServerPlaceholder: typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']
  NuxtLink: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']
  NuxtLoadingIndicator: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
  NuxtTime: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
  NuxtRouteAnnouncer: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
  NuxtAnnouncer: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-announcer")['default']
  NuxtImg: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
  NuxtPicture: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
  NuxtPage: typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']
  NoScript: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']
  Link: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']
  Base: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']
  Title: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']
  Meta: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']
  Style: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']
  Head: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']
  Html: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']
  Body: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']
  NuxtIsland: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']
  LazyAuthLogin: LazyComponent<typeof import("../../components/auth/Login.vue")['default']>
  LazyAuthLoginModal: LazyComponent<typeof import("../../components/auth/LoginModal.vue")['default']>
  LazyAuthRegister: LazyComponent<typeof import("../../components/auth/Register.vue")['default']>
  LazyCodeEditor: LazyComponent<typeof import("../../components/code/CodeEditor.vue")['default']>
  LazyCodePreview: LazyComponent<typeof import("../../components/code/CodePreview.vue")['default']>
  LazyCodeGeneratingOverlay: LazyComponent<typeof import("../../components/code/GeneratingOverlay.vue")['default']>
  LazyCodeMonacoEditorCore: LazyComponent<typeof import("../../components/code/MonacoEditorCore.vue")['default']>
  LazyCodeOptimizationPanel: LazyComponent<typeof import("../../components/code/OptimizationPanel.vue")['default']>
  LazyCodeSyntaxHighlighter: LazyComponent<typeof import("../../components/code/SyntaxHighlighter.vue")['default']>
  LazyCommonLoadingSpinner: LazyComponent<typeof import("../../components/common/LoadingSpinner.vue")['default']>
  LazyCommonModalDialog: LazyComponent<typeof import("../../components/common/ModalDialog.vue")['default']>
  LazyCommonToastNotification: LazyComponent<typeof import("../../components/common/ToastNotification.vue")['default']>
  LazyComparisonDiffViewer: LazyComponent<typeof import("../../components/comparison/DiffViewer.vue")['default']>
  LazyComparisonModelComparison: LazyComponent<typeof import("../../components/comparison/ModelComparison.vue")['default']>
  LazyComparisonVersionTimeline: LazyComponent<typeof import("../../components/comparison/VersionTimeline.vue")['default']>
  LazyEvaluationQualityGauge: LazyComponent<typeof import("../../components/evaluation/QualityGauge.vue")['default']>
  LazyEvaluationScoreCard: LazyComponent<typeof import("../../components/evaluation/ScoreCard.vue")['default']>
  LazyEvaluationSuggestionList: LazyComponent<typeof import("../../components/evaluation/SuggestionList.vue")['default']>
  LazyLayoutAppFooter: LazyComponent<typeof import("../../components/layout/AppFooter.vue")['default']>
  LazyLayoutAppHeader: LazyComponent<typeof import("../../components/layout/AppHeader.vue")['default']>
  LazyLayoutCardsArea: LazyComponent<typeof import("../../components/layout/CardsArea.vue")['default']>
  LazyLayoutHeroSection: LazyComponent<typeof import("../../components/layout/HeroSection.vue")['default']>
  LazyOptimizationCompareModal: LazyComponent<typeof import("../../components/optimization/CompareModal.vue")['default']>
  LazyOptimizationProgressModal: LazyComponent<typeof import("../../components/optimization/ProgressModal.vue")['default']>
  LazyUploadDragDropZone: LazyComponent<typeof import("../../components/upload/DragDropZone.vue")['default']>
  LazyUploadImagePreview: LazyComponent<typeof import("../../components/upload/ImagePreview.vue")['default']>
  LazyUploadImageUploader: LazyComponent<typeof import("../../components/upload/ImageUploader.vue")['default']>
  LazyNuxtWelcome: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
  LazyNuxtLayout: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
  LazyNuxtErrorBoundary: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
  LazyClientOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']>
  LazyDevOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']>
  LazyServerPlaceholder: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
  LazyNuxtLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
  LazyNuxtLoadingIndicator: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
  LazyNuxtTime: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
  LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
  LazyNuxtAnnouncer: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-announcer")['default']>
  LazyNuxtImg: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
  LazyNuxtPicture: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
  LazyNuxtPage: LazyComponent<typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']>
  LazyNoScript: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
  LazyLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']>
  LazyBase: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']>
  LazyTitle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']>
  LazyMeta: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']>
  LazyStyle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']>
  LazyHead: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']>
  LazyHtml: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']>
  LazyBody: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']>
  LazyNuxtIsland: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']>
}

declare module 'vue' {
  export interface GlobalComponents extends _GlobalComponents { }
}

export {}
