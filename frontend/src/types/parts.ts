/** Spare-parts inventory (API `/parts`, camelCase aligned with backend). */

export type PartUOM = 'unit' | 'liter';

export interface PartItem {
  id: string;
  reference: string;
  brand: string;
  name: string;
  barcode: string;
  quantity: number;
  uom: string;
  minimumQuantity?: number | null;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;
}

/** Body for POST /parts and PATCH /parts/:id (full field set on update, matching handler). */
export interface PartItemWriteBody {
  reference: string;
  brand: string;
  name: string;
  barcode: string;
  quantity: number;
  uom: PartUOM | string;
  minimumQuantity?: number | null;
}
