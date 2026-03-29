
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


export const AuthLogin: typeof import("../components/auth/Login.vue")['default']
export const AuthLoginModal: typeof import("../components/auth/LoginModal.vue")['default']
export const AuthRegister: typeof import("../components/auth/Register.vue")['default']
export const CodeEditor: typeof import("../components/code/CodeEditor.vue")['default']
export const CodePreview: typeof import("../components/code/CodePreview.vue")['default']
export const CodeGeneratingOverlay: typeof import("../components/code/GeneratingOverlay.vue")['default']
export const CodeMonacoEditorCore: typeof import("../components/code/MonacoEditorCore.vue")['default']
export const CodeOptimizationPanel: typeof import("../components/code/OptimizationPanel.vue")['default']
export const CodeSyntaxHighlighter: typeof import("../components/code/SyntaxHighlighter.vue")['default']
export const CommonLoadingSpinner: typeof import("../components/common/LoadingSpinner.vue")['default']
export const CommonModalDialog: typeof import("../components/common/ModalDialog.vue")['default']
export const CommonToastNotification: typeof import("../components/common/ToastNotification.vue")['default']
export const ComparisonDiffViewer: typeof import("../components/comparison/DiffViewer.vue")['default']
export const ComparisonModelComparison: typeof import("../components/comparison/ModelComparison.vue")['default']
export const ComparisonVersionTimeline: typeof import("../components/comparison/VersionTimeline.vue")['default']
export const EvaluationQualityGauge: typeof import("../components/evaluation/QualityGauge.vue")['default']
export const EvaluationScoreCard: typeof import("../components/evaluation/ScoreCard.vue")['default']
export const EvaluationSuggestionList: typeof import("../components/evaluation/SuggestionList.vue")['default']
export const LayoutAppFooter: typeof import("../components/layout/AppFooter.vue")['default']
export const LayoutAppHeader: typeof import("../components/layout/AppHeader.vue")['default']
export const LayoutCardsArea: typeof import("../components/layout/CardsArea.vue")['default']
export const LayoutHeroSection: typeof import("../components/layout/HeroSection.vue")['default']
export const OptimizationCompareModal: typeof import("../components/optimization/CompareModal.vue")['default']
export const OptimizationProgressModal: typeof import("../components/optimization/ProgressModal.vue")['default']
export const UploadDescEditorModal: typeof import("../components/upload/DescEditorModal.vue")['default']
export const UploadDragDropZone: typeof import("../components/upload/DragDropZone.vue")['default']
export const UploadPreviewImage: typeof import("../components/upload/PreviewImage.vue")['default']
export const UploadUploaderImage: typeof import("../components/upload/UploaderImage.vue")['default']
export const NuxtWelcome: typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']
export const NuxtLayout: typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
export const NuxtErrorBoundary: typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
export const ClientOnly: typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']
export const DevOnly: typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']
export const ServerPlaceholder: typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']
export const NuxtLink: typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']
export const NuxtLoadingIndicator: typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
export const NuxtTime: typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
export const NuxtRouteAnnouncer: typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
export const NuxtAnnouncer: typeof import("../node_modules/nuxt/dist/app/components/nuxt-announcer")['default']
export const NuxtImg: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
export const NuxtPicture: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
export const NuxtPage: typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']
export const NoScript: typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']
export const Link: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']
export const Base: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']
export const Title: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']
export const Meta: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']
export const Style: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']
export const Head: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']
export const Html: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']
export const Body: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']
export const NuxtIsland: typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']
export const LazyAuthLogin: LazyComponent<typeof import("../components/auth/Login.vue")['default']>
export const LazyAuthLoginModal: LazyComponent<typeof import("../components/auth/LoginModal.vue")['default']>
export const LazyAuthRegister: LazyComponent<typeof import("../components/auth/Register.vue")['default']>
export const LazyCodeEditor: LazyComponent<typeof import("../components/code/CodeEditor.vue")['default']>
export const LazyCodePreview: LazyComponent<typeof import("../components/code/CodePreview.vue")['default']>
export const LazyCodeGeneratingOverlay: LazyComponent<typeof import("../components/code/GeneratingOverlay.vue")['default']>
export const LazyCodeMonacoEditorCore: LazyComponent<typeof import("../components/code/MonacoEditorCore.vue")['default']>
export const LazyCodeOptimizationPanel: LazyComponent<typeof import("../components/code/OptimizationPanel.vue")['default']>
export const LazyCodeSyntaxHighlighter: LazyComponent<typeof import("../components/code/SyntaxHighlighter.vue")['default']>
export const LazyCommonLoadingSpinner: LazyComponent<typeof import("../components/common/LoadingSpinner.vue")['default']>
export const LazyCommonModalDialog: LazyComponent<typeof import("../components/common/ModalDialog.vue")['default']>
export const LazyCommonToastNotification: LazyComponent<typeof import("../components/common/ToastNotification.vue")['default']>
export const LazyComparisonDiffViewer: LazyComponent<typeof import("../components/comparison/DiffViewer.vue")['default']>
export const LazyComparisonModelComparison: LazyComponent<typeof import("../components/comparison/ModelComparison.vue")['default']>
export const LazyComparisonVersionTimeline: LazyComponent<typeof import("../components/comparison/VersionTimeline.vue")['default']>
export const LazyEvaluationQualityGauge: LazyComponent<typeof import("../components/evaluation/QualityGauge.vue")['default']>
export const LazyEvaluationScoreCard: LazyComponent<typeof import("../components/evaluation/ScoreCard.vue")['default']>
export const LazyEvaluationSuggestionList: LazyComponent<typeof import("../components/evaluation/SuggestionList.vue")['default']>
export const LazyLayoutAppFooter: LazyComponent<typeof import("../components/layout/AppFooter.vue")['default']>
export const LazyLayoutAppHeader: LazyComponent<typeof import("../components/layout/AppHeader.vue")['default']>
export const LazyLayoutCardsArea: LazyComponent<typeof import("../components/layout/CardsArea.vue")['default']>
export const LazyLayoutHeroSection: LazyComponent<typeof import("../components/layout/HeroSection.vue")['default']>
export const LazyOptimizationCompareModal: LazyComponent<typeof import("../components/optimization/CompareModal.vue")['default']>
export const LazyOptimizationProgressModal: LazyComponent<typeof import("../components/optimization/ProgressModal.vue")['default']>
export const LazyUploadDescEditorModal: LazyComponent<typeof import("../components/upload/DescEditorModal.vue")['default']>
export const LazyUploadDragDropZone: LazyComponent<typeof import("../components/upload/DragDropZone.vue")['default']>
export const LazyUploadPreviewImage: LazyComponent<typeof import("../components/upload/PreviewImage.vue")['default']>
export const LazyUploadUploaderImage: LazyComponent<typeof import("../components/upload/UploaderImage.vue")['default']>
export const LazyNuxtWelcome: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
export const LazyNuxtLayout: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
export const LazyNuxtErrorBoundary: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
export const LazyClientOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']>
export const LazyDevOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']>
export const LazyServerPlaceholder: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
export const LazyNuxtLink: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
export const LazyNuxtLoadingIndicator: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
export const LazyNuxtTime: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
export const LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
export const LazyNuxtAnnouncer: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-announcer")['default']>
export const LazyNuxtImg: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
export const LazyNuxtPicture: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
export const LazyNuxtPage: LazyComponent<typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']>
export const LazyNoScript: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
export const LazyLink: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']>
export const LazyBase: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']>
export const LazyTitle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']>
export const LazyMeta: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']>
export const LazyStyle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']>
export const LazyHead: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']>
export const LazyHtml: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']>
export const LazyBody: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']>
export const LazyNuxtIsland: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']>

export const componentNames: string[]
