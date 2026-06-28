export default function GlassBackground() {
  return (
    <div className="glass-bg" aria-hidden="true">
      <div className="glass-bg__surface" />
      <div className="glass-bg__sheen" />
      <div className="glass-bg__grid" />
      <div className="glass-bg__noise" />
    </div>
  );
}
