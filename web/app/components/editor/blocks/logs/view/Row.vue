<script lang="ts" setup>
import { cn } from "#layers/ui/app/lib/utils"
import HighlightedText from "./HighlightedText.vue"
import LevelPill from "./LevelPill.vue"
import RowDetails from "./RowDetails.vue"
import { summarizeLog } from "../rows"
import type { LineFormat, LogDialect, LogFilter } from "../dialect"
import { formatLogTime } from "../timestamp"
import { LogTimestampFormat } from "../utils"

const props = defineProps<{
	log: Log
	timestampFormat: LogTimestampFormat
	search: string
	dialect: LogDialect
	lineFilterFormats: LineFormat[]
	wrapLines?: boolean
}>()
const emit = defineEmits<{
	(e: "filter", filter: LogFilter): void
}>()
const expanded = defineModel<boolean>("expanded", { required: true })

const { d } = useI18n({ useScope: "global" })

const time = computed(() =>
	formatLogTime(
		props.log.time,
		expanded.value || props.timestampFormat === LogTimestampFormat.None
			? props.timestampFormat
			: LogTimestampFormat.NumericTime,
		d,
	),
)
const summary = computed(() => summarizeLog(props.log, props.dialect))
</script>
<template>
	<ShadcnUiCollapsible v-model:open="expanded">
		<ShadcnUiCollapsibleTrigger
			:class="
				cn(
					'group flex min-h-8 w-full cursor-pointer gap-2 px-2.5 text-left hover:bg-muted',
					props.wrapLines ? 'items-start py-1.5' : 'h-8 items-center',
					expanded && 'bg-muted',
				)
			"
		>
			<time
				v-if="time"
				:datetime="props.log.time"
				:title="props.log.time"
				class="max-w-[45%] shrink-0 truncate text-2sm whitespace-nowrap tabular-nums"
			>
				<span class="text-muted-foreground/70">{{ time.date }}</span>
				<template v-if="time.date">
					{{ " " }}
				</template>
				<span class="text-muted-foreground">{{ time.time }}</span>
			</time>
			<LevelPill
				:level="props.log.level"
				:class="props.wrapLines && 'mt-1.5'"
			/>
			<span
				:class="
					cn(
						'min-w-0 flex-1 text-2sm text-foreground',
						props.wrapLines
							? 'line-clamp-3 wrap-anywhere whitespace-pre-wrap'
							: 'truncate whitespace-nowrap',
					)
				"
			>
				<HighlightedText
					:text="summary.message || props.log.line"
					:search="props.search"
				/>
			</span>
			<!-- reserve the chevron so hovering cannot reflow a wrapped message -->
			<Icon
				name="mingcute:right-line"
				:class="
					cn(
						'-ml-2 h-3.5 w-0 shrink-0 text-muted-foreground/70 opacity-0 transition-transform duration-150 group-hover:ml-0 group-hover:w-3.5 group-hover:opacity-100',
						(expanded || props.wrapLines) && 'ml-0 w-3.5 opacity-100',
						expanded && 'rotate-90',
					)
				"
			/>
		</ShadcnUiCollapsibleTrigger>
		<ShadcnUiCollapsibleContent
			class="overflow-hidden data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down"
		>
			<RowDetails
				:log="props.log"
				:dialect="props.dialect"
				:line-filter-formats="props.lineFilterFormats"
				@filter="(filter) => emit('filter', filter)"
			/>
		</ShadcnUiCollapsibleContent>
	</ShadcnUiCollapsible>
</template>
