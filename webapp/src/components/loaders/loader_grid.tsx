import { grid } from 'ldrs'
grid.register()

export default function LoaderGrid() {
    return <l-grid
        size="75"
        speed="1.5"
        color="hsl(var(--primary))"
    ></l-grid>
}