import {defineStore} from 'pinia'

export class Mod {
    id: number
    gameId: number
    name: string
    description?: string

    constructor(partial: Partial<Mod>) {
        this.id = partial.id || 0;
        this.name = partial.name || '';
        this.gameId = partial.gameId || 0;
        this.description = partial.description || '';
    }
}

const mods: Mod[] = [
    ...Array.from({length: 100}, (_, i) => new Mod({
        id: i + 1,
        gameId: 1,
        name: `Mod ${Math.random().toString(36).substring(2, 8)}`,
        description: `Description ${Math.random().toString(36).substring(2, 15)}`
    }))
]

export const useModsStore = defineStore('mods', {
    state: () => ({
        mods: [] as Mod[],
        loading: false,
    }),
    actions: {
        async fetchMods(gameId: number, offset?:number, maxAmount?: number): Promise<Mod[]> {
            this.loading = true

            try {
                return new Promise<Mod[]>((resolve) => {
                    setTimeout(() => {
                        console.log(`Fetching mods for gameId: ${gameId}`)
                        const filteredMods = mods.filter(mod => mod.gameId === gameId)
                        const paginatedMods = offset !== undefined && maxAmount !== undefined
                            ? filteredMods.slice(offset, offset + maxAmount)
                            : filteredMods
                        this.mods = paginatedMods
                        resolve(paginatedMods)
                    }, 1000)
                })
            } catch (error) {
                console.error('Error fetching mods:', error)
                throw error
            } finally {
                this.loading = false
            }
        },
        getModCount(gameId: number): number {
            return mods.filter(mod => mod.gameId === gameId).length;
        }
    }
})
