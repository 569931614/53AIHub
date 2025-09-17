export const WEBSITE_STYLE = {
  WEBSITE: 'website',
  SOFTWARE: 'software'
} as const
export type WebsiteStyle = (typeof WEBSITE_STYLE)[keyof typeof WEBSITE_STYLE]

export const WEBSITE_STYLE_LABEL_MAP = new Map([
  [WEBSITE_STYLE.WEBSITE, 'guide.website'],
  [WEBSITE_STYLE.SOFTWARE, 'guide.software']
])
export const WEBSITE_STYLE_DEMO_MAP = new Map([
  [WEBSITE_STYLE.WEBSITE, '/images/template-website.png'],
  [WEBSITE_STYLE.SOFTWARE, '/images/template-software.png']
])
