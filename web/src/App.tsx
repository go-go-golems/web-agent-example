import { useCallback, useMemo, useState } from 'react';
import type { ComposerSlotProps, ChatWidgetRenderers } from '@pwchat/webchat';
import { ChatWidget } from '@pwchat/webchat';
import { ThinkingModeComposer } from './components/ThinkingModeComposer';
import { WebAgentThinkingModeCard } from './components/WebAgentThinkingModeCard';
import { DiscoDialogueCard } from './components/DiscoDialogueCard';

export function App() {
  const [mode, setMode] = useState('fast');
  const [discoEnabled, setDiscoEnabled] = useState(true);

  const buildOverrides = useCallback(() => {
    const middlewares: Array<{ name: string; config: Record<string, unknown> }> = [
      {
        name: 'webagent-thinking-mode',
        config: { mode },
      },
    ];
    if (discoEnabled) {
      middlewares.push({
        name: 'webagent-disco-dialogue',
        config: {},
      });
    }
    return { middlewares };
  }, [mode, discoEnabled]);

  const renderers: ChatWidgetRenderers = useMemo(
    () => ({
      webagent_thinking_mode: WebAgentThinkingModeCard,
      disco_dialogue_line: DiscoDialogueCard,
      disco_dialogue_check: DiscoDialogueCard,
      disco_dialogue_state: DiscoDialogueCard,
    }),
    [],
  );

  const components = useMemo(
    () => ({
      Composer: (props: ComposerSlotProps) => (
        <ThinkingModeComposer
          {...props}
          mode={mode}
          onModeChange={setMode}
          discoEnabled={discoEnabled}
          onDiscoEnabledChange={setDiscoEnabled}
        />
      ),
    }),
    [mode, discoEnabled],
  );

  return (
    <div className="app-shell">
      <div className="app-card">
        <ChatWidget renderers={renderers} components={components} buildOverrides={buildOverrides} />
      </div>
    </div>
  );
}
