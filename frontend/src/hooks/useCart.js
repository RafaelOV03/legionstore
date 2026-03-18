import { useState, useEffect } from 'react';
import { getCart } from '../services/api';

export const useCart = () => {
  const [cart, setCart] = useState({ items: [] });
  const [loading, setLoading] = useState(true);

  const refreshCart = async () => {
    try {
      setLoading(true);
      const data = await getCart();
      // Asegurar que items siempre sea un array
      setCart({
        ...data,
        items: data.items || []
      });
    } catch (error) {
      console.error('Error fetching cart:', error);
      setCart({ items: [] });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refreshCart();
  }, []);

  const cartItemsCount = cart.items?.reduce((sum, item) => sum + item.quantity, 0) || 0;
  
  const cartTotal = cart.items?.reduce(
    (sum, item) => sum + (item.product?.price || 0) * item.quantity,
    0
  ) || 0;

  return { cart, loading, refreshCart, cartItemsCount, cartTotal };
};
