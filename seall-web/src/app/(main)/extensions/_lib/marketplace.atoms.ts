import { atomWithStorage } from "jotai/utils"

// Default marketplace URL
export const DEFAULT_MARKETPLACE_URL = "https://raw.githubusercontent.com/Kolex06/Seall---based-on-seanime/main/marketplace.json"

// Atom to store the marketplace URL in localStorage
export const marketplaceUrlAtom = atomWithStorage<string>(
    "marketplace-url",
    DEFAULT_MARKETPLACE_URL,
    undefined,
    { getOnInit: true },
)
