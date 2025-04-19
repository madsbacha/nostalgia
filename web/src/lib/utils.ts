import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getInitials(name: string) {
  const filteredName = name.replace(/[^a-zA-Z0-9]/g, '');
  const uppercaseName = filteredName.toUpperCase();
  if (uppercaseName.length === 0) {
    return '';
  } else if (uppercaseName.length < 2) {
    return uppercaseName;
  }
  return uppercaseName.slice(0, 2);
}
