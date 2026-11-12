import { mountSuspended } from "@nuxt/test-utils/runtime"
import { describe, it } from "vitest"
import LevelPill from "./LevelPill.vue"
import { t } from "~/components/test-helpers"

const LEVELS = [
	LogLevel.Critical,
	LogLevel.Error,
	LogLevel.Warning,
	LogLevel.Info,
	LogLevel.Debug,
	LogLevel.Trace,
	LogLevel.Unknown,
]

describe("<LevelPill>", () => {
	it.for(LEVELS)("names the %s level", async (level, { expect }) => {
		const wrapper = await mountSuspended(LevelPill, { props: { level } })

		expect(wrapper.get(".sr-only").text()).toBe(
			t(`editor.logs.levels.${level}`),
		)
		expect(wrapper.get("span").attributes("title")).toBe(
			t(`editor.logs.levels.${level}`),
		)
	})

	it.for(LEVELS)(
		"abbreviates the %s level without duplicating its accessible name",
		async (level, { expect }) => {
			const wrapper = await mountSuspended(LevelPill, { props: { level } })
			const abbreviation = wrapper.get("[aria-hidden='true']")
			expect(abbreviation.text()).toBe(
				level === LogLevel.Unknown
					? "?"
					: t(`editor.logs.levels.${level}`).slice(0, 1),
			)
			expect(wrapper.get(".sr-only").text()).toBe(
				t(`editor.logs.levels.${level}`),
			)
		},
	)

	it("gives each level its own colour", async ({ expect }) => {
		const colours = await Promise.all(
			LEVELS.map(async (level) => {
				const wrapper = await mountSuspended(LevelPill, { props: { level } })

				return wrapper
					.get("span")
					.classes()
					.filter((c) => c.startsWith("text-") && c !== "text-2xs")
					.join(" ")
			}),
		)

		expect(new Set(colours).size).toBe(LEVELS.length)
	})
})
