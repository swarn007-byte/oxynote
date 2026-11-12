<script lang="ts" setup>
import { cn } from "#layers/ui/app/lib/utils"

const LEVEL_CLASSES: Record<LogLevel, string> = {
	[LogLevel.Critical]: "text-log-level-critical",
	[LogLevel.Error]: "text-log-level-error",
	[LogLevel.Warning]: "text-log-level-warning",
	[LogLevel.Info]: "text-log-level-info",
	[LogLevel.Debug]: "text-log-level-debug",
	[LogLevel.Trace]: "text-log-level-trace",
	[LogLevel.Unknown]: "text-log-level-unknown",
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
				'inline-flex size-2 shrink-0 rounded-full bg-current',
				LEVEL_CLASSES[props.level],
			)
		"
	>
		<span class="sr-only">{{ $t(`editor.logs.levels.${props.level}`) }}</span>
	</span>
</template>
