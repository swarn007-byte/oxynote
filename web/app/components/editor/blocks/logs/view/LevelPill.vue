<script lang="ts" setup>
import { cn } from "#layers/ui/app/lib/utils"

const LEVEL_CLASSES: Record<LogLevel, string> = {
	[LogLevel.Critical]: "bg-log-level-critical/10 text-log-level-critical",
	[LogLevel.Error]: "bg-log-level-error/10 text-log-level-error",
	[LogLevel.Warning]: "bg-log-level-warning/10 text-log-level-warning",
	[LogLevel.Info]: "bg-log-level-info/10 text-log-level-info",
	[LogLevel.Debug]: "bg-log-level-debug/10 text-log-level-debug",
	[LogLevel.Trace]: "bg-log-level-trace/10 text-log-level-trace",
	[LogLevel.Unknown]: "bg-log-level-unknown/10 text-log-level-unknown",
}

const props = defineProps<{
	level: LogLevel
}>()
</script>
<template>
	<span
		:title="$t(`editor.logs.levels.${props.level}`)"
		:class="
			cn(
				'inline-flex size-4.5 shrink-0 items-center justify-center rounded-full text-2xs font-semibold uppercase @sm:w-12.5',
				LEVEL_CLASSES[props.level],
			)
		"
	>
		<span aria-hidden="true" class="@sm:hidden">
			{{
				props.level === LogLevel.Unknown
					? "?"
					: $t(`editor.logs.levels.${props.level}`).slice(0, 1)
			}}
		</span>
		<span class="sr-only @sm:not-sr-only">
			{{ $t(`editor.logs.levels.${props.level}`) }}
		</span>
	</span>
</template>
