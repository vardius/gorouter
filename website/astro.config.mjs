// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
	site: 'https://vardius.github.io',
	base: '/gorouter',
	trailingSlash: 'always',
	// Pages import the runnable programs in ../examples, so the dev server may read outside this folder.
	vite: { server: { fs: { allow: ['..'] } } },
	integrations: [
		starlight({
			title: 'gorouter',
			description:
				'A fast HTTP router for Go, for net/http and fasthttp: zero-allocation matching, standard handlers and scoped middleware.',
			logo: { src: './src/assets/logo.png', alt: '' },
			favicon: '/favicon.ico',
			head: [
				{ tag: 'link', attrs: { rel: 'apple-touch-icon', href: '/gorouter/apple-touch-icon.png' } },
			],
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/vardius/gorouter' }],
			customCss: ['./src/styles/custom.css'],
			sidebar: [
				{
					label: 'Start here',
					items: ['getting-started', 'routing', 'middleware', 'mounting'],
				},
				{
					label: 'Guides',
					items: [
						'guides/static-files',
						'guides/authentication',
						'guides/error-handling',
						'guides/https',
						'guides/multiple-domains',
						'guides/fasthttp',
					],
				},
				{
					label: 'Reference',
					items: [
						'benchmarks',
						'migrating-from-v4',
						{ label: 'API reference', link: 'https://pkg.go.dev/github.com/vardius/gorouter/v5' },
					],
				},
			],
		}),
	],
});
