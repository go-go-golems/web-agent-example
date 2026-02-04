import { DefaultComposer, type ComposerSlotProps } from '@pwchat/webchat';

type Props = ComposerSlotProps & {
  mode: string;
  onModeChange: (next: string) => void;
};

export function ThinkingModeComposer({ mode, onModeChange, ...props }: Props) {
  return (
    <div className="webagent-composer">
      <div className="webagent-composer-toolbar">
        <label className="form-label" htmlFor="thinking-mode-select">
          Thinking mode
        </label>
        <select
          id="thinking-mode-select"
          className="form-select form-select-sm"
          value={mode}
          onChange={(e) => onModeChange(e.target.value)}
        >
          <option value="fast">Fast</option>
          <option value="deliberate">Deliberate</option>
          <option value="deep">Deep</option>
        </select>
      </div>
      <DefaultComposer {...props} />
    </div>
  );
}
